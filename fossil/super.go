package fossil

import (
	"sigint.ca/fs/internal/pack"
	"sigint.ca/fs/venti"
)

const (
	SuperMagic   = 0x2340a3b1
	SuperSize    = 512
	SuperVersion = 1
)

type Super struct {
	version   uint16
	epochLow  uint32
	epochHigh uint32
	qid       uint64      /* next qid */
	active    uint32      /* root of active file system */
	next      uint32      /* root of next snapshot to archive */
	current   uint32      /* root of snapshot currently archiving */
	last      venti.Score /* last snapshot successfully archived */
	name      [128]byte   /* label */
}

func SuperPack(s *Super, p []byte) {
	for i := 0; i < SuperSize; i++ {
		p[i] = 0
	}
	pack.U32PUT(p, SuperMagic)
	assert(s.version == SuperVersion)
	pack.U16PUT(p[4:], s.version)
	pack.U32PUT(p[6:], s.epochLow)
	pack.U32PUT(p[10:], s.epochHigh)
	pack.U64PUT(p[14:], s.qid)
	pack.U32PUT(p[22:], s.active)
	pack.U32PUT(p[26:], s.next)
	pack.U32PUT(p[30:], s.current)
	copy(p[34:], s.last[:venti.ScoreSize])
	copy(p[54:], s.name[:])
}

func SuperUnpack(s *Super, p []byte) error {
	*s = Super{}
	if pack.U32GET(p) != SuperMagic {
		goto Err
	}
	s.version = pack.U16GET(p[4:])
	if s.version != SuperVersion {
		goto Err
	}
	s.epochLow = pack.U32GET(p[6:])
	s.epochHigh = pack.U32GET(p[10:])
	s.qid = pack.U64GET(p[14:])
	if s.epochLow == 0 || s.epochLow > s.epochHigh || s.qid == 0 {
		goto Err
	}
	s.active = pack.U32GET(p[22:])
	s.next = pack.U32GET(p[26:])
	s.current = pack.U32GET(p[30:])
	copy(s.last[:], p[34:][:venti.ScoreSize])
	copy(s.name[:], p[54:])
	return nil

Err:
	*s = Super{}
	return EBadSuper
}
