/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package neocoin

import (
	"fmt"
	"github.com/shopspring/decimal"
	"testing"
)

func TestDecimalShit(t *testing.T) {
	num, _ := decimal.NewFromString("0.00005")
	num2 := num.Shift(-1)
	t.Logf("balance: %v\n", num2)
}

func orderHash(origins []string, addr string, start int) []string {
	fmt.Printf("find addr: %v\n", addr)
	fmt.Printf("origins: %v\n", origins)
	newHashs := make([]string, start)
	copy(newHashs, origins[:start])
	end := 0
	for i := start; i < len(origins); i++ {
		txAddr := origins[i]
		if txAddr == addr {
			newHashs = append(newHashs, txAddr)
			end = i
			break
		}
	}

	fmt.Printf("head: %v\n", newHashs)
	fmt.Printf("front: %v\n", origins[start:end])
	fmt.Printf("behind: %v\n", origins[end+1:])

	newHashs = append(newHashs, origins[start:end]...)
	newHashs = append(newHashs, origins[end+1:]...)
	return newHashs
}

func TestOrderHash(t *testing.T) {
	origins := []string{
		"c", "a", "b", "a", "b", "c", "b", "c", "a",
	}

	confused := []string{
		"b", "c", "a", "a", "c", "b", "b", "c", "a",
	}

	for i, w := range origins {
		confused = orderHash(confused, w, i)
	}

	fmt.Println(confused)
}

func TestGetTxId(t *testing.T) {
	str :="80000001c249bb0e8c4e02ed738eeafd5b61c180d9cb1d7633fcdb062aa3a01372ebf4050000019b7cffdaa674beae0f930ebe6085af9093e5fe56b34a5c220ccdcf6efc336fc500e1f50500000000205f46e5be17823bc84f060f545d55a56455f8790141407d27db1a9bbc6d7d156ad6d34b2499cdeba3515dcec7c38ad967bf164b0fe8e4948a828c140a7799317f0f1101022ea1ad9e4ccf2d731470be2413da72d6e05e232103df22a1f7263a5300ac68849696ab52ee79466de5c414e44fcc8ea43abd8dcb5fac"
	txId,err := GetTxId(str)
	if err!=nil{
		t.Error(err)
	}
	fmt.Println(txId)
}