package main

import (
	"context"
	"fmt"

	
	// "time"

	"log"
	"os"
	"strconv"
	"syscall"

	// "time"
	
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	// "github.com/hanwen/go-fuse/v2/fuse/pathfs"
)

// numberNode is a filesystem node representing an integer. Prime
// numbers are regular files, while composite numbers are directories
// containing all smaller numbers, eg.
//
//   $ ls -F  /tmp/x/6
//   2  3  4/  5
//
// the file system nodes are deduplicated using inode numbers. The
// number 2 appears in many directories, but it is actually the represented
// by the same numberNode{} object, with inode number 2.
//
//   $ ls -i1  /tmp/x/2  /tmp/x/8/6/4/2
//   2 /tmp/x/2
//   2 /tmp/x/8/6/4/2
//
type numberNode struct {
	// Must embed an Inode for the struct to work as a node.
	fs.Inode

	// num is the integer represented in this file/directory
	num int
}

// isPrime returns whether n is prime
func isPrime(n int) bool {
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func numberToMode(n int) uint32 {
	// prime numbers are files
	if isPrime(n) {
		return fuse.S_IFREG
	}
	// composite numbers are directories
	return fuse.S_IFDIR
}

// Ensure we are implementing the NodeReaddirer interface
var _ = (fs.NodeReaddirer)((*numberNode)(nil))

// Readdir is part of the NodeReaddirer interface
func (n *numberNode) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	r := make([]fuse.DirEntry, 0, n.num)
	for i := 2; i < n.num; i++ {
		d := fuse.DirEntry{
			Name: strconv.Itoa(i),
			Ino:  uint64(i),
			Mode: numberToMode(i),
		}
		fmt.Print("Inode :")
		fmt.Println(d.Ino)
		r = append(r, d)
	}
	fmt.Printf("path: %s \n",n.Inode.Path( n.Root()))
	// if(!n.Inode.IsRoot()){
	// 	fmt.Println("looping")
	// 	pinode := &n.Inode
	// 	pname := n.Inode.String()
	// 	for {
	// 		if pinode.IsRoot(){
	// 			fmt.Printf("break\n")
	// 			break
	// 		}
	// 		fmt.Printf("%s %s|||", time.Now(), pname)
	// 		pname, pinode = pinode.Parent()

	// 	}
	// 	fmt.Println("")
	// }
	

	return fs.NewListDirStream(r), 0
}

// Ensure we are implementing the NodeLookuper interface
var _ = (fs.NodeLookuper)((*numberNode)(nil))

// Lookup is part of the NodeLookuper interface
func (n *numberNode) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {

	i, err := strconv.Atoi(name)
	
	
	// fmt.Printf("%s isDir %t isRoot %t  looking up:  %s\n", time.Now().String(),n.Inode.IsDir(), n.Inode.IsRoot(), name)
	

	// if(!n.Inode.IsRoot()){
	// 	pname, _ := n.Inode.Parent()
	// 	fmt.Printf("Parent NodeName : %s\n", pname)
	// }
	// fmt.Println("---" + name)
	if err != nil {
		return nil, syscall.ENOENT
	}

	if i >= n.num || i <= 1 {
		return nil, syscall.ENOENT
	}

	stable := fs.StableAttr{
		Mode: numberToMode(i),
		// The child inode is identified by its Inode number.
		// If multiple concurrent lookups try to find the same
		// inode, they are deduplicated on this key.
		Ino: uint64(i),
	}
	operations := &numberNode{num: i}

	// The NewInode call wraps the `operations` object into an Inode.
	child := n.NewInode(ctx, operations, stable)
	// fmt.Println(n.Root().Parent())
	// In case of concurrent lookup requests, it can happen that operations !=
	// child.Operations().
	return child, 0
}

// ExampleDynamic is a whimsical example of a dynamically discovered
// file system.
func main() {
	// This is where we'll mount the FS
	mntDir := "/tmp/x"
	os.Mkdir(mntDir, 0755)
	root := &numberNode{num: 10}
	server, err := fs.Mount(mntDir, root, &fs.Options{
		MountOptions: fuse.MountOptions{
			// Set to true to see how the file system works.
			Debug: false,
		},
	})
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Mounted on %s", mntDir)
	log.Printf("Unmount by calling 'fusermount -u %s'", mntDir)

	// Wait until unmount before exiting
	server.Wait()
}
