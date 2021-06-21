package main

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

const (
	nonCollabD = 264770
	collabD    = 52770
)

var (
	collabs    = []string{"stakerA", "stakerB", "stakerC", "stakerD", "stakerE", "stakerF"}
	nonCollabs = []string{"stakerH", "stakerI", "stakerJ", "stakerK"}

	d1 = map[string]*farm{
		"stakerA": {Size: 100, Address: "FA2E8BC6av1qzqUkZBh2RZyt1J3heCSsXSXKQsK67wjXJwgr6a1E", Reward: 132.65},
		"stakerB": {Size: 120, Address: "FA3GvuGLCN6EPUFd7uWG4eARJKEa9jF4Eiwg43iBHRKM4m4H5p63", Reward: 159.18},
		"stakerC": {Size: 200, Address: "FA2T81BybsfD411Su2FsMmRHHVkKoWLWAg7HdhQABm3YEx2qvL5L", Reward: 265.30},
		"stakerD": {Size: 50, Address: "FA3dGXc5MFp9b745G1xSP5gw3TGhtSk5DbJTBMEYUd79g6gLpfSJ", Reward: 66.33},
		"stakerE": {Size: 300, Address: "FA2DDF8jfC5oEicNevDknFdEAMC6FVjgp2nDejahhYKUKucn6JM2", Reward: 397.95},
		"stakerF": {Size: 2000, Address: "FA3p9DvQq4oaCp6x9ep93TEaJJ74PwHdSHETXr6AsKDm3iwtQHhi", Reward: 2653.02},
		"stakerG": {Size: 50000, Address: "FA3pztZxdBy1fQHjLbKh3MpU9PNJUhvcYCcRZAp5vUPtnJeUaphr", Reward: 13951.35},
	}

	d2 = map[string]*farm{
		"stakerH": {Size: 45000, Address: "FA2JZZXcdgBEpbqdmVNsmvSSfratu3k5sJoEAcxAEj3UPov37KCY", Reward: 11897.12},
		"stakerI": {Size: 55000, Address: "FA2SKfVuwarw4tt9ZSWDkfUrGqpwnpVMUVN4EbnkmTvYdXt3YMWN", Reward: 14540.92},
		"stakerJ": {Size: 60000, Address: "FA2RAhx382pNa3RSqFnRyRFKGrkrHQtovEQpF4D5EJCzhcSanwUK", Reward: 15862.82},
		"stakerK": {Size: 52000, Address: "FA2zktSMfomTGbqkZWUDTgtTeA51E65XEGz9ymQtwEM2k1SmfWrW", Reward: 13747.78},
	}

	delegatorsMap = make(map[string]*farm)
)

func init() {
	// build signatures before running tests
	err := buildSignatures(d1, collabs, "collabs")
	if err != nil {
		panic(err)
	}

	err = buildSignatures(d2, nonCollabs, "nonCollabs")
	if err != nil {
		panic(err)
	}
	for _, v := range collabs {
		delegatorsMap[v] = d1[v]
	}

	// create sig for G
	delegatrosBytes, _ := json.Marshal(delegatorsMap)
	sig, err := signData(d1["stakerG"].Address, delegatrosBytes)
	if err != nil {
		panic(fmt.Sprintf("[ERROR] unable to sign data, err : %v", err))
	}
	d1["stakerG"].Content = string(delegatrosBytes)
	d1["stakerG"].Sig, d1["stakerG"].PubKey = sig.Signature, sig.PubKey

}

func TestCalculateReward(t *testing.T) {

	for k, v := range d1 {
		c := collabD
		if k == "stakerG" {
			c = nonCollabD
		}

		res, err := calculateReward(d1, k, c, v.Size)
		if err != nil {
			log.Printf("[ERROR] unable to calculate reward, delegator : %v , err : %v", k, err)
		}
		if fmt.Sprintf("%.2f", res) != fmt.Sprintf("%.2f", v.Reward) {
			t.Errorf("reward not correct for staker : %v , stakerReward : %v , calculatedReward : %v", k, v.Reward, res)
		}
	}

	for k, v := range d2 {
		res, err := calculateReward(d2, k, nonCollabD, v.Size)
		if err != nil {
			log.Printf("[ERROR] unable to calculate reward, non delegator : %v , err : %v", k, err)
		}

		if fmt.Sprintf("%.2f", res) != fmt.Sprintf("%.2f", v.Reward) {
			t.Errorf("reward not correct for staker : %v , stakerReward : %v , calculatedReward : %v", k, v.Reward, res)
		}
	}

}
