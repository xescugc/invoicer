package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/filesystem"
)

func initializeFilesystemBilling() (billing.Billing, error) {
	baseDir, err := xdg.DataFile("invoicer/invoicer.data")
	if err != nil {
		panic(err)
	}

	// The lib (xdg) requires the .data file (or a file) at the end
	// so to get the correct path we need to call dir onto that as
	// the previous baseDir has the 'invoice.data' at the end
	baseDir = filepath.Dir(baseDir)

	ur := filesystem.NewUserRepository(baseDir)
	cr, err := filesystem.NewCustomerRepository(baseDir)
	if err != nil {
		return nil, err
	}

	ir, err := filesystem.NewInvoiceRepository(baseDir)
	if err != nil {
		return nil, err
	}

	tr, err := filesystem.NewTemplateRepository(baseDir)
	if err != nil {
		return nil, err
	}

	return billing.New(ur, cr, ir, tr), nil
}

// DefaultEditor is vim because we're adults ;)
const DefaultEditor = "vim"

// openFileInEditor opens filename in a text editor.
func openFileInEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = DefaultEditor
	}

	// Get the full executable path for the editor.
	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	cmd := exec.Command(executable, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// CaptureInputFromEditor opens a temporary file in a text editor and returns
// the written bytes on success or an error on failure. It handles deletion
// of the temporary file behind the scenes.
// The data will be marshaled into the file to fill
func captureInputFromEditor(data interface{}) ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return []byte{}, err
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return []byte{}, nil
	}

	file.Write(b)

	filename := file.Name()

	// Defer removal of the temporary file in case any of the next steps fail.
	defer os.Remove(filename)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if err = openFileInEditor(filename); err != nil {
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
