package dataset

import (
	"compress/bzip2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const datasetExt = ".bz2"

const datasetBaseURL = "https://dumps.wikimedia.org/enwiki/20231020/"

// ErrNotFound is returned by [Load] if the requested dataset is not on disk.
var ErrNotFound = errors.New("dataset: not found")

// Names contains the name of known dataset.
var Names = []string{
	"enwiki-20231020-pages-articles1.xml-p1p41242",
	"enwiki-20231020-pages-articles2.xml-p41243p151573",
	"enwiki-20231020-pages-articles3.xml-p151574p311329",
	"enwiki-20231020-pages-articles4.xml-p311330p558391",
	"enwiki-20231020-pages-articles5.xml-p558392p958045",
	"enwiki-20231020-pages-articles6.xml-p958046p1483661",
}

// Dataset represents a group of wiki pages.
type Dataset struct {
	name string
	size int64
	f    *os.File
	io.Reader
}

// Load loads a dataset from disk.
func Load(name string) (*Dataset, error) {
	var err error
	var f *os.File
	for _, path := range []string{name, name + datasetExt} {
		f, err = os.Open(path)
		if err == nil {
			break
		}
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		return nil, fmt.Errorf("open %v: %v", path, err)
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	}

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat: %v", err)
	}

	return &Dataset{
		name:   name,
		size:   info.Size(),
		f:      f,
		Reader: bzip2.NewReader(f),
	}, nil
}

// ProgressFunc is an optional function enabling the caller to wrap the reader
// and follow the download progress.
type ProgressFunc func(n int64, r io.Reader) io.Reader

// Download downloads the given dataset.
func Download(name string, progress ProgressFunc) (*Dataset, error) {
	url := datasetBaseURL + name + datasetExt

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get %v: %v", url, err)
	}
	defer resp.Body.Close()

	f, err := os.Create(name + datasetExt)
	if err != nil {
		return nil, fmt.Errorf("open %v: %v", name, err)
	}

	var reader io.Reader = resp.Body
	if progress != nil {
		reader = progress(resp.ContentLength, reader)
	}
	if _, err := io.Copy(f, reader); err != nil {
		return nil, fmt.Errorf("copy: %v", err)
	}

	if _, err := f.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("seek: %v", err)
	}

	return &Dataset{
		name:   name,
		size:   resp.ContentLength,
		f:      f,
		Reader: bzip2.NewReader(f),
	}, nil
}

// Name returns the name of the dataset.
func (d *Dataset) Name() string { return d.name }

// Size returns the size of the dataset.
func (d *Dataset) Size() int64 { return d.size }

// Close closes the dataset.
func (d *Dataset) Close() error {
	f := d.f
	d.f = nil
	d.Reader = nil
	return f.Close()
}

// Datasets represents a group of dataset.
type Datasets []*Dataset

// Close closes all the datasets.
func (datasets Datasets) Close() error {
	var errs []error
	for _, d := range datasets {
		if err := d.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
