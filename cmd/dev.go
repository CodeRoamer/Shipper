package cmd

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"io"

	"github.com/codegangsta/cli"

	"github.com/coderoamer/shipper/modules/log"
	"github.com/coderoamer/shipper/modules/setting"
)

var CmdDev = cli.Command{
	Name:  "dev",
	Usage: "Download all the assets and put them in the right places",
	Description: `Prepare everything you need to start your shipping dev`,
	Action: runDev,
	Flags:  []cli.Flag{},
}

func runDev(*cli.Context) {
	setting.NewConfigContext()
	log.Trace("Log path: %s", setting.LogRootPath)

	log.Trace("\n1. Use bower to download all the assets...\n")
	cmd := exec.Command("bower", "install")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Trace("\n2. Merge Assets into Public Folder...\n")

	// TODO: Check if bower_components folder exist!
	bowerDir := path.Join(setting.StaticRootPath, "bower_components")
	publicDir := path.Join(setting.StaticRootPath, "public")

	jsVendorDir := path.Join(publicDir, "js/vendor")
	fontsDir := path.Join(publicDir, "fonts")
	cssDir := path.Join(publicDir, "css")

	err = os.Mkdir(jsVendorDir, os.ModePerm)
	if err != nil && os.IsNotExist(err){
		log.Fatal("Error Make Folder public/js/vendor: %s", err.Error())
	}
	err = os.Mkdir(fontsDir, os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		log.Fatal("Error Make Folder public/fonts: %s", err.Error())
	}
	err = os.Mkdir(cssDir, os.ModePerm)
	if err != nil && os.IsNotExist(err) {
		log.Fatal("Error Make Folder public/css: %s", err.Error())
	}

	file, err := os.Create(filepath.Join(jsVendorDir, "script.js"))
	defer file.Close()

	if err != nil {
		log.Fatal("Error Creating script.js in public/js/vendor: %s", err.Error())
	}

	foundationFile, err := os.Create(filepath.Join(jsVendorDir, "foundation.js"))
	defer foundationFile.Close()

	if err != nil {
		log.Fatal("Error Creating foundation.js in public/js/vendor: %s", err.Error())
	}

	// jquery
	writeToFile(filepath.Join(bowerDir, "jquery/dist","jquery.js"), file)
	// other assets
	writeToFile(filepath.Join(bowerDir, "nprogress","nprogress.js"), file)
	writeToFile(filepath.Join(bowerDir, "loader.js","loader.js"), file)
	writeToFile(filepath.Join(bowerDir, "keymaster","keymaster.js"), file)
	writeToFile(filepath.Join(bowerDir, "jquery-file-upload/js","jquery.fileupload.js"), file)
	writeToFile(filepath.Join(bowerDir, "ic-ajax/dist/globals","main.js"), file)
	writeToFile(filepath.Join(bowerDir, "codemirror/lib","codemirror.js"), file)
	writeToFile(filepath.Join(bowerDir, "codemirror/addon/mode","overlay.js"), file)
	writeToFile(filepath.Join(bowerDir, "codemirror/mode/markdown","markdown.js"), file)
	writeToFile(filepath.Join(bowerDir, "codemirror/mode/gfm","gfm.js"), file)
	writeToFile(filepath.Join(bowerDir, "showdown/src","showdown.js"), file)
	writeToFile(filepath.Join(bowerDir, "Countable","Countable.js"), file)
	// ember and functional framework
	writeToFile(filepath.Join(bowerDir, "validator-js","validator.js"), file)
	writeToFile(filepath.Join(bowerDir, "lodash/dist","lodash.underscore.js"), file)
	writeToFile(filepath.Join(bowerDir, "handlebars","handlebars.js"), file)
	writeToFile(filepath.Join(bowerDir, "ember","ember.js"), file)
	writeToFile(filepath.Join(bowerDir, "ember-data","ember-data.js"), file)
	writeToFile(filepath.Join(bowerDir, "ember-resolver/dist","ember-resolver.js"), file)
	writeToFile(filepath.Join(bowerDir, "ember-load-initializers","ember-load-initializers.js"), file)
	// foundation framework
	writeToFile(filepath.Join(bowerDir, "fastclick/lib","fastclick.js"), foundationFile)
	writeToFile(filepath.Join(bowerDir, "jquery-placeholder","jquery.placeholder.js"), foundationFile)
	writeToFile(filepath.Join(bowerDir, "jquery.cookie","jquery.cookie.js"), foundationFile)
	writeToFile(filepath.Join(bowerDir, "foundation/js","foundation.js"), foundationFile)

	copyToFolder(filepath.Join(bowerDir, "nprogress", "nprogress.css"), filepath.Join(cssDir, "nprogress.css"))
	copyToFolder(filepath.Join(bowerDir, "fontawesome/css", "font-awesome.css"), filepath.Join(cssDir, "font-awesome.css"))
	copyToFolder(filepath.Join(bowerDir, "foundation/css", "foundation.css"), filepath.Join(cssDir, "foundation.css"))
	copyToFolder(filepath.Join(bowerDir, "foundation/css", "foundation.css.map"), filepath.Join(cssDir, "foundation.css.map"))
	copyToFolder(filepath.Join(bowerDir, "foundation/css", "normalize.css"), filepath.Join(cssDir, "normalize.css"))
	copyToFolder(filepath.Join(bowerDir, "foundation/css", "normalize.css.map"), filepath.Join(cssDir, "normalize.css.map"))
	copyToFolder(filepath.Join(bowerDir, "fontawesome/fonts", "FontAwesome.otf"), filepath.Join(fontsDir, "FontAwesome.otf"))
	copyToFolder(filepath.Join(bowerDir, "fontawesome/fonts", "fontawesome-webfont.eot"), filepath.Join(fontsDir, "fontawesome-webfont.eot"))
	copyToFolder(filepath.Join(bowerDir, "fontawesome/fonts", "fontawesome-webfont.svg"), filepath.Join(fontsDir, "fontawesome-webfont.svg"))
	copyToFolder(filepath.Join(bowerDir, "fontawesome/fonts", "fontawesome-webfont.ttf"), filepath.Join(fontsDir, "fontawesome-webfont.ttf"))
	copyToFolder(filepath.Join(bowerDir, "fontawesome/fonts", "fontawesome-webfont.woff"), filepath.Join(fontsDir, "fontawesome-webfont.woff"))
	copyToFolder(filepath.Join(bowerDir, "modernizr","modernizr.js"), filepath.Join(jsVendorDir, "modernizr.js"))


	log.Trace("\nDev Config Complete!")
}

func copyToFolder(oldPath, newPath string) {
	fcopy, err := os.Create(newPath)

	defer fcopy.Close()
	if err != nil {
		log.Fatal("Error Copy %s to %s: %s", oldPath, newPath, err.Error())
	}

	writeToFile(oldPath, fcopy)

	log.Info("finish Copy %s to %s...", filepath.Base(oldPath), filepath.Base(newPath))
}

func writeToFile(finPath string, fout *os.File) {
	fin, err := os.Open(finPath)

	buf := make([]byte, 1024)

	if err != nil {
		log.Fatal("Error Reading %s: %s", finPath, err.Error())
	}

	for {
		n, err := fin.Read(buf)
		if n == 0 && err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Error Reading %s: %s", finPath, err.Error())
		}
		_, err = fout.Write(buf[:n])
		if err != nil {
			log.Fatal("Error Reading %s: %s", finPath, err.Error())
		}
	}

	log.Info("finish reading %s...", filepath.Base(finPath))
}

