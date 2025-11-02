package capaFachadaServices

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"log"

	"github.com/tcolgate/mp3"
)

// GuardarArchivoCancion copia un archivo de audio desde srcPath a la carpeta "canciones/"
// y lo renombra usando el título (sanitizado) como "titulo.mp3". Devuelve la ruta destino.
func GuardarArchivoCancion(srcPath, titulo string) (string, error) {
	if srcPath == "" || titulo == "" {
		return "", nil
	}

	// Asegurar carpeta
	destDir := "canciones"
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}

	// Comprobar título para nombre de archivo
	re := regexp.MustCompile(`[^\p{L}0-9 '\-_]+`)
	name := re.ReplaceAllString(titulo, "_")
	// limpiar guiones bajos repetidos y trim espacios al inicio / fin
	name = strings.TrimSpace(name)
	name = regexp.MustCompile(`_+`).ReplaceAllString(name, "_")
	dest := filepath.Join(destDir, name+".mp3")

	// Copiar
	in, err := os.Open(srcPath)
	if err != nil {
		log.Printf("Error abriendo archivo origen %s: %v", srcPath, err)
		return "", err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		log.Printf("Error creando archivo destino %s: %v", dest, err)
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		log.Printf("Error copiando archivo a %s: %v", dest, err)
		return "", err
	}

	return dest, nil
}

// Guarda el contenido del reader en la carpeta canciones/ con nombre tratado
func GuardarArchivoDesdeReader(r io.Reader, titulo string) (string, error) {
	if r == nil || titulo == "" {
		return "", nil
	}
	destDir := "canciones"
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", err
	}
	re := regexp.MustCompile(`[^\p{L}0-9 '\-_]+`)
	name := re.ReplaceAllString(titulo, "_")
	name = strings.TrimSpace(name)
	name = regexp.MustCompile(`_+`).ReplaceAllString(name, "_")
	dest := filepath.Join(destDir, name+".mp3")

	out, err := os.Create(dest)
	if err != nil {
		log.Printf("Error creando archivo destino %s: %v", dest, err)
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		log.Printf("Error guardando archivo desde reader a %s: %v", dest, err)
		return "", err
	}
	return dest, nil
}

// ObtenerDuracionFFprobe intenta extraer la duración del archivo usando ffprobe (ffmpeg).
// Devuelve la duración en formato mm:ss
func ObtenerDuracionFFprobe(path string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", path)
	out, err := cmd.CombinedOutput()
	s := strings.TrimSpace(string(out))
	if err != nil {
		log.Printf("ffprobe error: %s", s)
		log.Printf("ffprobe exec error: %v — attempting Go-based fallback", err)
		if d, ferr := obtenerDuracionMP3(path); ferr == nil && d != "" {
			return d, nil
		} else if ferr != nil {
			log.Printf("fallback mp3 parse fallo: %v", ferr)
		}
		return "", fmt.Errorf("ffprobe fallo: %w -- output: %s", err, s)
	}
	if s == "" {
		log.Printf("ffprobe retorno vacio para %s, intentando fallback", path)
		if d, ferr := obtenerDuracionMP3(path); ferr == nil && d != "" {
			return d, nil
		} else if ferr != nil {
			log.Printf("fallback mp3 parse fallo %v", ferr)
		}
		return "", nil
	}
	secs, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return "", err
	}
	// formato mm:ss (redondeo al segundo)
	total := int(secs + 0.5)
	mins := total / 60
	sec := total % 60
	return strconv.Itoa(mins) + ":" + pad2(sec), nil
}

func pad2(v int) string {
	if v < 10 {
		return "0" + strconv.Itoa(v)
	}
	return strconv.Itoa(v)
}

// Intenta extraer la duración leyendo los frames MP3 usando un decoder en Go.
func obtenerDuracionMP3(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	d := mp3.NewDecoder(f)
	var frame mp3.Frame
	var total time.Duration
	for {
		if err := d.Decode(&frame, nil); err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		total += frame.Duration()
	}
	if total == 0 {
		return "", nil
	}
	secs := int(total.Seconds() + 0.5)
	mins := secs / 60
	sec := secs % 60
	return strconv.Itoa(mins) + ":" + pad2(sec), nil
}
