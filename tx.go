// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btcutil

import (
	"bytes"
	"io"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
)

// TxIndexUnknown is the value returned for a transaction index that is unknown.
// This is typically because the transaction has not been inserted into a block
// yet.
const TxIndexUnknown = -1

// Tx defines a bitcoin transaction that provides easier and more efficient
// manipulation of raw transactions.  It also memoizes the hash for the
// transaction on its first access so subsequent accesses don't have to repeat
// the relatively expensive hashing operations.
type Tx struct {
	msgTx         *wire.MsgTx     // Underlying MsgTx
	txHash        *chainhash.Hash // Cached transaction hash
	txHashWitness *chainhash.Hash // Cached transaction witness hash
	txHasWitness  *bool           // If the transaction has witness data
	txIndex       int             // Position within a block or TxIndexUnknown
}

type TxNew struct {
	msgTxNew      *wire.MsgTxNew     // Underlying MsgTx
	txHash        *chainhash.Hash // Cached transaction hash
	txHashWitness *chainhash.Hash // Cached transaction witness hash
	txHasWitness  *bool           // If the transaction has witness data
	txIndex       int             // Position within a block or TxIndexUnknown
}

// MsgTx returns the underlying wire.MsgTx for the transaction.
func (t *Tx) MsgTx() *wire.MsgTx {
	// Return the cached transaction.
	return t.msgTx
}

func (t *TxNew) MsgTxNew() *wire.MsgTxNew {
	// Return the cached transaction.
	return t.msgTxNew
}

func (t *TxNew) MsgTx() *wire.MsgTx {
	// Return the cached transaction.
	return t.msgTxNew.CreateMsgTx()
}

// Hash returns the hash of the transaction.  This is equivalent to
// calling TxHash on the underlying wire.MsgTx, however it caches the
// result so subsequent calls are more efficient.
func (t *Tx) Hash() *chainhash.Hash {
	// Return the cached hash if it has already been generated.
	if t.txHash != nil {
		return t.txHash
	}

	// Cache the hash and return it.
	hash := t.msgTx.TxHash()
	t.txHash = &hash
	return &hash
}

func (t *TxNew) Hash() *chainhash.Hash {
	// Return the cached hash if it has already been generated.
	if t.txHash != nil {
		return t.txHash
	}

	// Cache the hash and return it.
	hash := t.msgTxNew.TxHash()
	t.txHash = &hash
	return &hash
}

// WitnessHash returns the witness hash (wtxid) of the transaction.  This is
// equivalent to calling WitnessHash on the underlying wire.MsgTx, however it
// caches the result so subsequent calls are more efficient.
func (t *Tx) WitnessHash() *chainhash.Hash {
	// Return the cached hash if it has already been generated.
	if t.txHashWitness != nil {
		return t.txHashWitness
	}

	// Cache the hash and return it.
	hash := t.msgTx.WitnessHash()
	t.txHashWitness = &hash
	return &hash
}

// HasWitness returns false if none of the inputs within the transaction
// contain witness data, true false otherwise. This equivalent to calling
// HasWitness on the underlying wire.MsgTx, however it caches the result so
// subsequent calls are more efficient.
func (t *Tx) HasWitness() bool {
	if t.txHashWitness != nil {
		return *t.txHasWitness
	}

	hasWitness := t.msgTx.HasWitness()
	t.txHasWitness = &hasWitness
	return hasWitness
}

func (t *TxNew) HasWitness() bool {
	if t.txHasWitness != nil {
		return *t.txHasWitness
	}

	hasWitness := t.msgTxNew.HasWitness()
	t.txHasWitness = &hasWitness
	return hasWitness
}

// Index returns the saved index of the transaction within a block.  This value
// will be TxIndexUnknown if it hasn't already explicitly been set.
func (t *Tx) Index() int {
	return t.txIndex
}

// Index returns the saved index of the transaction within a block.  This value
// will be TxIndexUnknown if it hasn't already explicitly been set.
func (t *TxNew) Index() int {
	return t.txIndex
}

// SetIndex sets the index of the transaction in within a block.
func (t *Tx) SetIndex(index int) {
	t.txIndex = index
}

// SetIndex sets the index of the transaction in within a block.
func (t *TxNew) SetIndex(index int) {
	t.txIndex = index
}

// NewTx returns a new instance of a bitcoin transaction given an underlying
// wire.MsgTx.  See Tx.
func NewTx(msgTx *wire.MsgTx) *Tx {
	return &Tx{
		msgTx:   msgTx,
		txIndex: TxIndexUnknown,
	}
}

func NewTxNew(msgTx *wire.MsgTx) *Tx {
	return &Tx{
		msgTx:   msgTx,
		txIndex: 	TxIndexUnknown,
	}
}

// NewTxFromBytes returns a new instance of a bitcoin transaction given the
// serialized bytes.  See Tx.
func NewTxFromBytes(serializedTx []byte) (*TxNew, error) {
	br := bytes.NewReader(serializedTx)
	return NewTxFromReader(br)
}

// NewTxFromReader returns a new instance of a bitcoin transaction given a
// Reader to deserialize the transaction.  See Tx.
func NewTxFromReader(r io.Reader) (*TxNew, error) {
	// Deserialize the bytes into a MsgTx.
	var msgTxNew wire.MsgTxNew
	err := msgTxNew.Deserialize(r)
	if err != nil {
		return nil, err
	}

	t := TxNew{
		msgTxNew:   &msgTxNew,
		txIndex: TxIndexUnknown,
	}
	return &t, nil
}
