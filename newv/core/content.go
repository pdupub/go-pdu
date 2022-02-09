// Copyright 2021 The PDU Authors
// This file is part of the PDU library.
//
// The PDU library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PDU library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PDU library. If not, see <http://www.gnu.org/licenses/>.

package core

const (
	QCFmtStringTEXT       = 1
	QCFmtStringURL        = 2
	QCFmtStringJSON       = 3
	QCFmtStringInt        = 4
	QCFmtStringFloat      = 5
	QCFmtStringHexAddress = 6

	QCFmtBytesSignature = 33

	QCFmtImagePNG = 65
	QCFmtImageJPG = 66
	QCFmtImageBMP = 67

	QCFmtAudioWAV = 97
	QCFmtAudioMP3 = 98

	QCFmtVideoMP4 = 129
)

// QContent is one piece of data in Quantum,
// all variables should be in alphabetical order.
type QContent struct {
	Data   []byte `json:"data,omitempty"`
	Format int    `json:"fmt"`
}

func NewContent(fmt int, data []byte) (*QContent, error) {
	return &QContent{Format: fmt, Data: data}, nil
}
