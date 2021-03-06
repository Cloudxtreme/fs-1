package main

import "testing"

func BenchmarkFileWrite(b *testing.B) {
	if err := testAllocFsys(); err != nil {
		b.Fatalf("testAllocFsys: %v", err)
	}
	defer testCleanupFsys()

	fsys, err := getFsys("testfs")
	if err != nil {
		b.Errorf("get fsys: %v", err)
		return
	}
	fs := fsys.getFs()
	fsys.put()

	buf := make([]byte, 10*8192)
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i % 256)
	}

	dir, err := fs.openFile("/active")
	if err != nil {
		b.Errorf("open dir: %v", err)
		return
	}
	file, err := dir.create("writetest", 0644, "none")
	dir.decRef()
	if err != nil {
		b.Errorf("create: %v", err)
		return
	}

	for i := 0; i < b.N; i++ {
		if _, err := file.write(buf, len(buf), 0, "none"); err != nil {
			b.Errorf("write: %v", err)
			break
		}
	}

	if err := file.remove("none"); err != nil {
		b.Errorf("remove: %v", err)
		return
	}
	file.decRef()
}

func BenchmarkFileRead(b *testing.B) {
	if err := testAllocFsys(); err != nil {
		b.Fatalf("testAllocFsys: %v", err)
	}
	defer testCleanupFsys()

	fsys, err := getFsys("testfs")
	if err != nil {
		b.Errorf("get fsys: %v", err)
		return
	}
	fs := fsys.getFs()
	fsys.put()

	buf := make([]byte, 10*8192)
	for i := 0; i < len(buf); i++ {
		buf[i] = byte(i % 256)
	}

	dir, err := fs.openFile("/active")
	if err != nil {
		b.Errorf("open dir: %v", err)
		return
	}
	file, err := dir.create("readtest", 0644, "none")
	dir.decRef()
	if err != nil {
		b.Errorf("create: %v", err)
		return
	}
	if _, err := file.write(buf, len(buf), 0, "none"); err != nil {
		b.Errorf("write: %v", err)
		return
	}

	memset(buf, 0)

	for i := 0; i < b.N; i++ {
		n, err := file.read(buf, 0)
		if err != nil {
			b.Errorf("read: %v", err)
			break
		}
		if n != len(buf) {
			b.Errorf("bad read length: got %d, want %d", n, len(buf))
			break
		}
	}

	if err := file.remove("none"); err != nil {
		b.Errorf("remove: %v", err)
		return
	}
	file.decRef()
}
func BenchmarkFileCreateDelete(b *testing.B) {
	if err := testAllocFsys(); err != nil {
		b.Fatalf("testAllocFsys: %v", err)
	}
	defer testCleanupFsys()

	fsys, err := getFsys("testfs")
	if err != nil {
		b.Errorf("get fsys: %v", err)
		return
	}
	fs := fsys.getFs()
	fsys.put()

	dir, err := fs.openFile("/active")
	if err != nil {
		b.Errorf("open dir: %v", err)
		return
	}

	for i := 0; i < b.N; i++ {
		file, err := dir.create("createtest", 0644, "none")
		if err != nil {
			b.Errorf("create: %v", err)
			return
		}
		if err := file.remove("none"); err != nil {
			b.Errorf("remove: %v", err)
			return
		}
		file.decRef()
	}
	dir.decRef()
}
