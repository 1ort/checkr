package mmdb

import (
    "archive/tar"
    "compress/gzip"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
    "time"

    "github.com/oschwald/maxminddb-golang"
)

type MaxMindClient struct {
	licenseKey    string
	HttpClient    *http.Client
	Verbose       bool
	directoryPath string
}

func NewDownloader(licenseKey string, directoryPath string) *MaxMindClient {
	return &MaxMindClient{
		licenseKey:    licenseKey,
		HttpClient:    http.DefaultClient,
		directoryPath: directoryPath,
	}
}

func OpenLocal(directoryPath, editionID string) (*maxminddb.Reader, error) {
    dbPath := filepath.Join(directoryPath, editionID+".mmdb")
    db, err := maxminddb.Open(dbPath)
    if err != nil{
        return nil, err
    }
    return db, nil
}

func (mm *MaxMindClient) Latest(editionID string) (*maxminddb.Reader, error) {
	mm.log("Opening local database")
	dbPath := filepath.Join(mm.directoryPath, editionID+".mmdb")
	err := os.MkdirAll(mm.directoryPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	db, err := maxminddb.Open(dbPath)
	if os.IsNotExist(err) {
		mm.log("File not found, downloading database")
		return mm.downloadAndOpen(editionID)
	}
	mm.log("Opened")
	if err != nil {
		db.Close()
		return nil, err
	}
	mm.log("Checking if the database is up to date")
	up2date, err := mm.IsUpToDate(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("check db is up-to-date: %w", err)
	}
	if !up2date {
		db.Close()
		mm.log("Database is outdated, updating")
		return mm.downloadAndOpen(editionID)
	}
	mm.log("Database is up to date")
	return db, nil
}

func (mm *MaxMindClient) IsUpToDate(db *maxminddb.Reader) (bool, error) {
	editionID := db.Metadata.DatabaseType
	locallastModified, err := mm.getLocalLastModified(editionID)
	if os.IsNotExist(err) {
		mm.log("Local Last-modified.txt file not found")
		return false, nil
	}
	if err != nil {
		return false, err
	}
	mm.log("Local build time: " + locallastModified.Format("2006-01-02 15:04"))
	lastModified, err := mm.GetLastModified(editionID)
	if err != nil {
		return false, fmt.Errorf("retrive Last-modified value: %w", err)
	}
	mm.log("Latest build time: " + lastModified.Format("2006-01-02 15:04"))
	return !locallastModified.Before(lastModified), nil

}

func (mm *MaxMindClient) GetLastModified(editionID string) (time.Time, error) {
	downloadUrl := BuildDownloadUrl(mm.licenseKey, editionID)
	res, err := mm.HttpClient.Head(downloadUrl)
	if err != nil {
		return time.Time{}, fmt.Errorf("send HEAD request: %w", err)
	}
	lastModified := res.Header.Get("Last-Modified")
	return parselastModified(lastModified)

}

func (mm *MaxMindClient) getLocalLastModified(editionID string) (time.Time, error) {
	txtPath := filepath.Join(mm.directoryPath, editionID+".Last-Modified.txt")
	content, err := os.ReadFile(txtPath)
	if err != nil {
		return time.Time{}, err
	}
	return parselastModified(string(content))
}

func (mm *MaxMindClient) saveLocalLastModified(editionID, lastModified string) error {
	txtPath := filepath.Join(mm.directoryPath, editionID+".Last-Modified.txt")
	f, err := os.Create(txtPath)
	if err != nil {
		return fmt.Errorf("create .Last-Modified.txt file: %v", err)
	}
	defer f.Close()
	_, wrerr := f.WriteString(lastModified)
	if wrerr != nil {
		return fmt.Errorf("write .Last-Modified.txt file: %v", err)
	}
	return nil
}

func BuildDownloadUrl(licenseKey, editionID string) string {
	base, err := url.Parse("https://download.maxmind.com/app/geoip_download")
	if err != nil {
		return ""
	}
	params := url.Values{}
	params.Add("edition_id", editionID)
	params.Add("license_key", licenseKey)
	params.Add("suffix", "tar.gz")
	base.RawQuery = params.Encode()

	return base.String()
}

func GetBuildTime(db *maxminddb.Reader) time.Time {
	buildEpoch := int64(db.Metadata.BuildEpoch)
	buildTime := time.Unix(buildEpoch, 0)
	return buildTime
}

func parselastModified(lastModified string) (time.Time, error) {
	if lastModified == "" {
		return time.Time{}, errors.New("no Last-Modified header found")
	}

	t, err := time.ParseInLocation(time.RFC1123, lastModified, time.UTC)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing time: %w", err)
	}

	return t, nil
}

func (mm *MaxMindClient) downloadAndOpen(editionID string) (*maxminddb.Reader, error) {
	url := BuildDownloadUrl(mm.licenseKey, editionID)
	destPath := filepath.Join(mm.directoryPath, editionID+".mmdb")
	lastModified, err := downloadAndUnpack(mm.HttpClient, url, destPath)
	if err != nil {
		return nil, fmt.Errorf("download new database: %w", err)
	}
	err = mm.saveLocalLastModified(editionID, lastModified)
	if err != nil {
		mm.log("error writing Last-modified.txt file: " + fmt.Sprint(err))
	}
	db, err := maxminddb.Open(destPath)
	if err != nil {
		return nil, fmt.Errorf("open new database: %w", err)
	}

	return db, nil
}

func (mm *MaxMindClient) log(message string) {
	if mm.Verbose {
		log.Println(message)
	}
}

func downloadAndUnpack(client *http.Client, url, destPath string) (string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("request file %v: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		st := resp.Status
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("request file %v: %v: %v", url, st, string(b))
	}
	gzr, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gzip.NewReader: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("read archive: %w", err)
		}
		if filepath.Ext(header.Name) == ".mmdb" {
			f, err := os.Create(destPath)
			if err != nil {
				return "", fmt.Errorf("create file: %v", err)
			}
			defer f.Close()
			if _, err := io.Copy(f, tr); err != nil {
				return "", fmt.Errorf("copy file from archive: %w", err)
			}
			lastModified := resp.Header.Get("Last-Modified")
			return lastModified, nil
		}
	}
	return "", errors.New(".mmdb file not found in archive")
}