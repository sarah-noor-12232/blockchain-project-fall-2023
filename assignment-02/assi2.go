package main

import (
    "fmt"
    "sync"
)

type Transaction struct {
    ID string
}

type Block struct {
    Transactions []Transaction
}

type Node struct {
    ID            string
    Transactions  []Transaction
    Blocks        []Block
    Neighbors     []*Node
    Network       *Network
}

type Network struct {
    Nodes []*Node
}

func NewNode(id string, network *Network) *Node {
    node := &Node{
        ID:      id,
        Network: network,
    }
    network.Nodes = append(network.Nodes, node)
    return node
}

func (node *Node) AddNeighbor(neighbor *Node) {
    node.Neighbors = append(node.Neighbors, neighbor)
}

func (node *Node) BroadcastTransaction(transaction Transaction) {
    node.Transactions = append(node.Transactions, transaction)
    for _, neighbor := range node.Neighbors {
        neighbor.ReceiveTransaction(transaction)
    }
}

func (node *Node) ReceiveTransaction(transaction Transaction) {
    node.Transactions = append(node.Transactions, transaction)
}

func (node *Node) MineBlock() {
    block := Block{
        Transactions: node.Transactions,
    }
    node.Blocks = append(node.Blocks, block)
    node.Transactions = []Transaction{}
    for _, neighbor := range node.Neighbors {
        neighbor.ReceiveBlock(block)
    }
}

func (node *Node) ReceiveBlock(block Block) {
    node.Blocks = append(node.Blocks, block)
}

func main() {
    network := &Network{}
    nodes := make([]*Node, 10)
    for i := range nodes {
        nodes[i] = NewNode(fmt.Sprintf("Node %d", i+1), network)
        if i > 0 {
            nodes[i].AddNeighbor(nodes[i-1])
        }
    }
    nodes[0].BroadcastTransaction(Transaction{ID: "Transaction 1"})
    nodes[0].MineBlock()
    for _, node := range nodes {
        fmt.Printf("%s has %d transactions and %d blocks\n", node.ID, len(node.Transactions), len(node.Blocks))
    }
}