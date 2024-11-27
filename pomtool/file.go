package pomtool

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
)

type File string

func (f File) String() string {
	return string(f)
}

func (f File) Child(name string) File {
	path := f.String() + "/" + name
	return File(path)
}

func (f File) elements() []string {
	const (
		c1 = "\\"
		c2 = "/"
	)
	str := f.String()
	str = strings.ReplaceAll(str, c1, c2)
	return strings.Split(str, c2)
}

func (f File) stringify(elist []string) string {
	b := &strings.Builder{}
	s0 := ""
	sep := ""
	if strings.HasPrefix(f.String(), "/") {
		// unix-like
		s0 = "/"
		sep = "/"
	} else {
		// win32
		s0 = ""
		sep = "\\"
	}
	for idx, item := range elist {
		if idx == 0 {
			b.WriteString(s0)
		} else {
			b.WriteString(sep)
		}
		b.WriteString(item)
	}
	return b.String()
}

func (f File) normalize(src []string) []string {
	dst := make([]string, 0)
	for _, item := range src {
		if item == "" {
			continue
		} else if item == "." {
			continue
		} else if item == ".." {
			size := len(dst)
			if size > 0 {
				dst = dst[0 : size-1]
			} else {
				path := f.stringify(src)
				panic("cannot normalize path: " + path)
			}
		} else {
			dst = append(dst, item)
		}
	}
	return dst
}

func (f File) Normalize() File {
	elist := f.elements()
	elist = f.normalize(elist)
	str := f.stringify(elist)
	return File(str)
}

func (f File) Name() string {
	list := f.elements()
	size := len(list)
	for i := size - 1; i >= 0; i-- {
		item := strings.TrimSpace(list[i])
		if item != "" {
			return item
		}
	}
	return ""
}

func (f File) Parent() File {
	f2 := f + "/.."
	return f2.Normalize()
}

func (f File) Exists() bool {
	_, err := os.Stat(f.String())
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}

func (f File) IsDir() bool {
	st, err := os.Stat(f.String())
	if err == nil {
		return st.IsDir()
	}
	return false
}

func (f File) IsFile() bool {
	st, err := os.Stat(f.String())
	if err == nil {
		return !st.IsDir()
	}
	return false
}

func (f File) ReadBinary() ([]byte, error) {
	return os.ReadFile(f.String())
}

func (f File) CopyTo(dst File) error {

	src := f
	flag := os.O_WRONLY | os.O_TRUNC
	perm := fs.ModePerm

	if dst.IsDir() {
		return fmt.Errorf("dest path is a dir : %s", dst.String())
	}

	i, err := os.Open(src.String())
	if err != nil {
		return err
	}
	defer i.Close()

	o := i
	o = nil
	if dst.Exists() {
		o, err = os.OpenFile(dst.String(), flag, perm)
	} else {
		o, err = os.Create(dst.String())
	}
	if err != nil {
		return err
	}
	defer o.Close()

	_, err = io.Copy(o, i)
	return err
}

func (f File) Mkdirs() error {
	if f.Exists() {
		return nil
	}
	perm := fs.ModePerm
	path := f.String()
	return os.MkdirAll(path, perm)
}
