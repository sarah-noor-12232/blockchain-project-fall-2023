package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type MerkelNode struct {
	left        *MerkelNode
	right       *MerkelNode
	transaction string
}

type MerkelTree struct {
	Root *MerkelNode
}

func AddNode(transac string, arr *[10]MerkelTree, merkelindex *int) {

	var ind int
	ind = *merkelindex
	var tree *MerkelNode = arr[ind].Root

	var newitem MerkelNode

	newitem.transaction = transac
	newitem.right = nil
	newitem.left = nil

	if tree == nil {
		fmt.Println("Inserting nil ", transac, "\n")
		arr[ind].Root = &newitem
		ind++

	} else {

		fmt.Println("Inserting ", transac, "\n")
		InsertItem(tree, &newitem)
		ind++
	}

	*merkelindex = ind
}

func InsertItem(tree *MerkelNode, item *MerkelNode) {

	if len(item.transaction) <= len(tree.transaction) {

		if tree.left == nil {
			tree.left = item
			return
		} else {
			InsertItem(tree.left, item)
			return
		}
	} else if len(item.transaction) > len(tree.transaction) {

		if tree.right == nil {
			tree.right = item
			return
		} else {
			InsertItem(tree.right, item)
			return
		}

	}
}

func DisplayTree(tree *MerkelNode) {

	if tree == nil {
		return
	}
	fmt.Println(tree.transaction, "\n")

	if tree.left != nil {
		DisplayTree(tree.left)
	}

	if tree.right != nil {
		DisplayTree(tree.right)
	}
}

func Update2(tree *MerkelNode, prev string, now string) {

	if tree == nil {
		return
	}
	if prev == tree.transaction {
		tree.transaction = now
	}

	if tree.left != nil {
		Update2(tree.left, prev, now)
	}

	if tree.right != nil {
		Update2(tree.right, prev, now)
	}

}

func Update(arr *[10]MerkelTree, merkelindex *int, prev string, now string) {

	var ind int
	ind = *merkelindex
	var i = 0
	for i = 0; i < ind+1; i++ {
		var tree *MerkelNode = arr[i].Root

		Update2(tree, prev, now)

	}

	*merkelindex = ind

}

func Traversal2(tree *MerkelNode, tt *string) {

	if tree == nil {
		return
	}
	*tt += tree.transaction

	if tree.left != nil {
		Traversal2(tree.left, tt)
	}

	if tree.right != nil {
		Traversal2(tree.right, tt)
	}

}

func Traversal(arr *[10]MerkelTree, merkelindex *int) string {

	var ind int
	ind = *merkelindex
	var alltransactions string
	var i = 0
	for i = 0; i < ind+1; i++ {
		var tree *MerkelNode = arr[i].Root

		Traversal2(tree, &alltransactions)

	}

	*merkelindex = ind

	return alltransactions

}

func DisplayMerkelTree(arr *[10]MerkelTree, merkelindex *int) {
	var ind int
	ind = *merkelindex
	var i = 0

	fmt.Println("------------------Merkel Tree elements-------------------------------------")
	for i = 0; i < ind+1; i++ {
		var tree *MerkelNode = arr[i].Root

		if tree == nil {
			fmt.Println("arr[", i, "] has no elements\n")

		} else {
			fmt.Println("arr[", i, "] has elements\n")
			DisplayTree(tree)
		}
	}

	*merkelindex = ind

}

func CreateMerkelTree(arr *[10]MerkelTree, merkelindex *int) {

	for i := 0; i < 10; i++ {
		arr[i].Root = nil
	}

	for i := 0; i < nooftransactions; i++ {
		AddNode(transactionslist[i], arr, merkelindex)
	}

}

type block struct {
	arr           [10]MerkelTree
	merkelindex   int
	id            int
	nonce         string
	previous_hash string
	current_hash  string
}

type blockchain struct {
	list []*block
}

func NewBlock(x int) *block {
	//fmt.Println("------------------------fdsfdsfdsfds-----------------------------------")
	tempblock := new(block)
	tempblock.id = x
	tempblock.nonce = "0"
	tempblock.merkelindex = 0
	CreateMerkelTree(&tempblock.arr, &tempblock.merkelindex)
	return tempblock
}

func VerifyChain(chain *blockchain) bool {
	var temp = ""
	var check = true
	for i := 0; i < len(chain.list); i++ {
		tt := Traversal(&chain.list[i].arr, &chain.list[i].merkelindex)

		var attributes string
		attributes += strconv.Itoa(chain.list[i].id)
		attributes += tt + chain.list[i].previous_hash
		total_sum := sha256.Sum256([]byte(attributes))
		temp = fmt.Sprintf("%x", total_sum)

		if temp != chain.list[i].current_hash {
			check = false
			fmt.Printf("Previous block has been tampered, i.e. Block # %d\n", i)
			break

		}
	}

	if check == false {
		fmt.Println("error occured")
	} else {
		fmt.Printf("Blocks verified. No tampering\n")
	}
	return check
}

func VerifyBlock(myblock *block, NewBlock *block) bool {

	if NewBlock.current_hash == myblock.current_hash && NewBlock.nonce == myblock.nonce {
		return true
	}

	return false
}
func MineBlock(blocklist *blockchain) {

	for j := 0; j < len(blocklist.list); j++ {
		print("to match:", blocklist.list[j].current_hash, "\n")
		for i := 0; ; i++ {
			temp := sha256.Sum256([]byte(strconv.Itoa(i)))
			noncex := fmt.Sprintf("%x", temp)
			dum := noncex[:3]
			fmt.Println("dum:", dum)
			fmt.Println(strings.Contains(blocklist.list[j].current_hash, dum))

			if strings.Contains(blocklist.list[j].current_hash, dum) == true {
				blocklist.list[j].nonce = dum
				break

			}

		}

	}

}

func MineSingleBlock(NewBlock *block) {

	print("to match:", NewBlock.current_hash, "\n")
	for i := 0; ; i++ {
		temp := sha256.Sum256([]byte(strconv.Itoa(i)))
		noncex := fmt.Sprintf("%x", temp)
		dum := noncex[:3]
		fmt.Println("dum:", dum)
		fmt.Println(strings.Contains(NewBlock.current_hash, dum))

		if strings.Contains(NewBlock.current_hash, dum) == true {
			NewBlock.nonce = dum
			fmt.Println("founddd")
			break

		}

	}

}
func CalculateSingleBlockHash(NewBlock *block) {

	tt := Traversal(&NewBlock.arr, &NewBlock.merkelindex)
	var attributes string
	attributes += strconv.Itoa(NewBlock.id)
	attributes += tt + NewBlock.previous_hash
	total_sum := sha256.Sum256([]byte(attributes))
	NewBlock.current_hash = fmt.Sprintf("%x", total_sum)
}

func CalculateHash(chain *blockchain) {

	for i := 0; i < len(chain.list); i++ {
		tt := Traversal(&chain.list[i].arr, &chain.list[i].merkelindex)
		var attributes string
		attributes += strconv.Itoa(chain.list[i].id)
		attributes += tt + chain.list[i].previous_hash
		total_sum := sha256.Sum256([]byte(attributes))
		chain.list[i].current_hash = fmt.Sprintf("%x", total_sum) // formating to string
		if i < len(chain.list)-1 {
			chain.list[i+1].previous_hash = fmt.Sprintf("%x", total_sum) //storing current block hash to next block in its previous hash var
		}

	}
}

func (blocklist *blockchain) addblock(x int) *block {
	tempblock := NewBlock(x)

	if VerifyChain(blocklist) {
		blocklist.list = append(blocklist.list, tempblock)
		CalculateHash(blocklist)

		fmt.Printf("block addition in chain successful\n")
	} else {
		fmt.Printf(" error. block addition unsuccessful.\n")
		return nil
	}
	return tempblock
}

func DisplayBlocks(blocklist *blockchain) {
	fmt.Println("")

	for i := 0; i < len(blocklist.list); i++ {
		fmt.Printf("Block id:%d\n\n", blocklist.list[i].id)
		DisplayMerkelTree(&blocklist.list[i].arr, &blocklist.list[i].merkelindex)
		fmt.Println("nonce value : \n", blocklist.list[i].nonce)
		fmt.Println("current hash: \n", blocklist.list[i].current_hash)
		fmt.Println("previous hash: \n", blocklist.list[i].previous_hash)

	}

	fmt.Println("")

}

func ChangeBlock(chain *blockchain, x int) { // updating on basis of id value as identifier

	found := false
	for i := 0; i < len(chain.list); i++ {

		if x == chain.list[i].id {

			var now string
			var prev string
			fmt.Println("Enter transaction to change\n")

			fmt.Scanln(&prev)

			fmt.Println("Enter Updated value\n")

			fmt.Scanln(&now)
			fmt.Println("Updated successfully\n")
			Update(&chain.list[i].arr, &chain.list[i].merkelindex, prev, now)
			found = true
		}
	}
	if found == false {
		fmt.Println("error. Couldnt Update. block not found")
	}
	return
}

func contains(sarr []string, str string) bool {
	for _, i := range sarr {
		if i == str {
			return true
		}
	}
	return false
}

func dealconn(values ...interface{}) {
	for _, val := range values {
		_ = val
	}
}

var transactionslist []string
var nooftransactions int
var nodes []string
var check bool
var done bool

func main() {

	nooftransactions = 0
	check = false
	done = false
	conn, err := net.Dial("tcp", "localhost:8001")
	if err != nil && check == false {
		// error handling
	}
	rcvdarr := make([]byte, 200)
	conn.Read(rcvdarr)
	var rcvbuffer string = string(rcvdarr)
	var rcvdmesssage string = strings.Split(rcvbuffer, "\000")[0]

	if rcvdmesssage == "none" {

		ln, err := net.Listen("tcp", conn.LocalAddr().String())
		if err != nil {
			log.Fatal(err)
		}
		for {
			conn2, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			fmt.Println("acception succesful for node 1")
			go handleConnection(conn2)
			fmt.Println("transaction list:")
			for i := 0; i < nooftransactions; i++ {
				fmt.Println(transactionslist[i])
			}

		}

	} else {

		//Dial
		fmt.Println("receivedbuff:", rcvbuffer)
		nodes = strings.Split(rcvbuffer, ",")
		for i := 1; i <= len(nodes); i++ {
			var ip string = strings.Split(nodes[i-1], ":")[0]
			var portstr string = strings.Split(nodes[i-1], ":")[1]
			fmt.Println(ip, portstr)
			var port string = ""
			for k := 0; k < len(portstr); k++ {
				var j int = 0
				if portstr[k] >= '0' && portstr[k] <= '9' {
					port += string(portstr[k])
					j++
				} else {
					break
				}
			}

			var str string = "localhost:" + port
			connection, err := net.Dial("tcp", str)
			dealconn(connection)
			if err != nil {
				fmt.Println("Error at node:", i)
				fmt.Println("url:", ip, port)

			}
			fmt.Println("Connected to:", ip, port)
			msg := "hello from new node"
			connection.Write([]byte(msg))

		}

		ln, err := net.Listen("tcp", conn.LocalAddr().String())
		if err != nil {
			fmt.Println("listen error")
			log.Fatal(err)
		}
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err)
				continue

			}
			fmt.Println("acception succesful")
			go handleConnection(conn)
			fmt.Println("transaction list:")
			for i := 0; i < nooftransactions; i++ {
				fmt.Println(transactionslist[i])
			}

		}

	}
}

func handleConnection(conn net.Conn) {

	rcvdarr := make([]byte, 500)

	if rcvdarr != nil && check == false {

		conn.Read(rcvdarr)

		var rcvbuffer string = string(rcvdarr)
		var rcvdmesssage string = strings.Split(rcvbuffer, "\000")[0]
		fmt.Println("message received:", rcvdmesssage)

		if rcvdmesssage[0] == 't' {
			fmt.Println("yes sir")
			if !contains(transactionslist, rcvdmesssage) {
				transactionslist = append(transactionslist, rcvdmesssage)
				nooftransactions++

			} else {
				fmt.Println(rcvdmesssage + " already present here\n")
			}

		}
		if nooftransactions != 0 && nooftransactions <= 5 && check == false {

			fmt.Println("flooding:")

			for i := 1; i <= len(nodes); i++ {
				var ip string = strings.Split(nodes[i-1], ":")[0]
				var portstr string = strings.Split(nodes[i-1], ":")[1]
				fmt.Println(ip, portstr)
				var port string = ""
				for k := 0; k < len(portstr); k++ {
					var j int = 0
					if portstr[k] >= '0' && portstr[k] <= '9' {
						port += string(portstr[k])
						j++
					} else {
						break
					}
				}

				var str string = "localhost:" + port
				connection, err := net.Dial("tcp", str)
				dealconn(connection)
				if err != nil {
					fmt.Println("Error at node:", i)
					fmt.Println("url:", ip, port)

				}
				fmt.Println("Connected to:", ip, port)
				msg := transactionslist[nooftransactions-1]
				fmt.Println("send msg:", msg)
				connection.Write([]byte(msg))

			}
			for i := 0; i < nooftransactions; i++ {
				fmt.Println(transactionslist[i])
			}

			if nooftransactions == 5 {
				check = true
			}
		}

	}

	if check == true && done == false {

		blockid := 13
		nodeblock := NewBlock(blockid)
		CalculateSingleBlockHash(nodeblock)

		MineSingleBlock(nodeblock)

		fmt.Println("sending mined block:")

		for i := 1; i <= len(nodes); i++ {
			var ip string = strings.Split(nodes[i-1], ":")[0]
			var portstr string = strings.Split(nodes[i-1], ":")[1]
			fmt.Println(ip, portstr)
			var port string = ""
			for k := 0; k < len(portstr); k++ {
				var j int = 0
				if portstr[k] >= '0' && portstr[k] <= '9' {
					port += string(portstr[k])
					j++
				} else {
					break
				}
			}

			var str string = "localhost:" + port
			connection2, err := net.Dial("tcp", str)
			dealconn(connection2)
			if err != nil {
				fmt.Println("Error at node:", i)
				fmt.Println("url:", ip, port)

			}
			fmt.Println("Connected to:", ip, port)
			connection2.Write([]byte("bandar"))

			//bin_buf := new(bytes.Buffer)
			// create a encoder object
			//gobobj := gob.NewEncoder(bin_buf)
			//encode buffer and marshal it into a gob object
			//fmt.Println("noooooooooooooooooooooo  ", nodeblock.nonce)
			//gobobj.Encode(nodeblock)

			//connection.Write(bin_buf.Bytes())

		}

		done = true

	}

	if check == true && done == true {

		conn.Read(rcvdarr)

		fmt.Println("haaaaaaaaallloooooo:", string(rcvdarr))

		//tempbuff := bytes.NewBuffer(rcvdarr)

		//NewBlock := NewBlock(13)

		//gobobj := gob.NewDecoder(tempbuff)

		//gobobj.Decode(NewBlock)

		//fmt.Println("received block nonce:", NewBlock.nonce)

		//var rcvbuffer string = string(rcvdarr)
		//var rcvdmesssage string = strings.Split(rcvbuffer, "\000")[0]
		//fmt.Println("received after mining:", rcvdmesssage)
	}

}
