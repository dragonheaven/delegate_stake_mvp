package main

import (
	"encoding/json"
	"errors"
	"fmt"
	fac "github.com/FactomProject/factom"
	facWallet "github.com/FactomProject/factom/wallet"
	"github.com/FactomProject/factomd/common/primitives"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	delegators = make(map[string]*farm)

	collaborator = "stakerG"

	collaborators    = []string{"stakerA", "stakerB", "stakerC", "stakerD", "stakerE", "stakerF"}
	nonCollaborators = []string{"stakerH", "stakerI", "stakerJ", "stakerK"}

	reward      = 70000
	maxFarmSize = 40000
)

type (
	farm struct {
		Size    int     `json:"size"`
		Address string  `json:"address"`
		Content string  `json:"content"`
		Sig     []byte  `json:"sig"`
		PubKey  []byte  `json:"pub_key"`
		Reward  float64 `json:"reward"`
	}
)

func main() {

	// input data
	farms := map[string]*farm{
		"stakerA": {Size: 100, Address: "FA2E8BC6av1qzqUkZBh2RZyt1J3heCSsXSXKQsK67wjXJwgr6a1E"},
		"stakerB": {Size: 120, Address: "FA3GvuGLCN6EPUFd7uWG4eARJKEa9jF4Eiwg43iBHRKM4m4H5p63"},
		"stakerC": {Size: 200, Address: "FA2T81BybsfD411Su2FsMmRHHVkKoWLWAg7HdhQABm3YEx2qvL5L"},
		"stakerD": {Size: 50, Address: "FA3dGXc5MFp9b745G1xSP5gw3TGhtSk5DbJTBMEYUd79g6gLpfSJ"},
		"stakerE": {Size: 300, Address: "FA2DDF8jfC5oEicNevDknFdEAMC6FVjgp2nDejahhYKUKucn6JM2"},
		"stakerF": {Size: 2000, Address: "FA3p9DvQq4oaCp6x9ep93TEaJJ74PwHdSHETXr6AsKDm3iwtQHhi"},
		"stakerG": {Size: 50000, Address: "FA3pztZxdBy1fQHjLbKh3MpU9PNJUhvcYCcRZAp5vUPtnJeUaphr"},
		"stakerH": {Size: 45000, Address: "FA2JZZXcdgBEpbqdmVNsmvSSfratu3k5sJoEAcxAEj3UPov37KCY"},
		"stakerI": {Size: 55000, Address: "FA2SKfVuwarw4tt9ZSWDkfUrGqpwnpVMUVN4EbnkmTvYdXt3YMWN"},
		"stakerJ": {Size: 60000, Address: "FA2RAhx382pNa3RSqFnRyRFKGrkrHQtovEQpF4D5EJCzhcSanwUK"},
		"stakerK": {Size: 52000, Address: "FA2zktSMfomTGbqkZWUDTgtTeA51E65XEGz9ymQtwEM2k1SmfWrW"},
	}

	err := buildSignatures(farms, nonCollaborators, "nonCollabs")
	if err != nil {
		panic(err)
	}

	err = buildSignatures(farms, collaborators, "collabs")
	if err != nil {
		panic(err)
	}

	// build data for G
	for _, v := range collaborators {
		delegators[v] = farms[v]
	}

	// create sig for G
	delegatrosBytes, _ := json.Marshal(delegators)
	sig, err := signData(farms[collaborator].Address, delegatrosBytes)
	if err != nil {
		panic(fmt.Sprintf("[ERROR] unable to sign data, address : %v , err : %v", farms[collaborator].Address, err))
	}
	farms[collaborator].Content = string(delegatrosBytes)
	farms[collaborator].Sig, farms[collaborator].PubKey = sig.Signature, sig.PubKey

	collabDivisor := farms[collaborator].Size

	nonCollabDivisor := 0
	topStakers := map[string]int{}

	for k, v := range farms {
		if v.Size >= maxFarmSize {
			topStakers[k] = v.Size
			nonCollabDivisor += v.Size
		}
	}

	stakersReward := map[string]float64{}

	for k, v := range topStakers {
		stakersReward[k], err = calculateReward(farms, k, nonCollabDivisor, v)
		if err != nil {
			log.Printf("[ERROR] unable to calculate reward, farmer : %v , err : %v", k, err)
		}
	}

	for _, collab := range collaborators {
		stakersReward[collab], err = calculateReward(farms, collab, collabDivisor, farms[collab].Size)
		if err != nil {
			log.Printf("[ERROR] unable to calculate reward, delegator : %v , err : %v", collab, err)
		}
	}

	for k, v := range stakersReward {
		fmt.Printf("%v) %.2f$ Reward\n", k, v)
	}
}

func buildSignatures(farms map[string]*farm, data []string, typ string) error {
	var (
		staker string
	)

	for _, v := range data {

		switch typ {
		case "nonCollabs":
			staker = v
		case "collabs":
			staker = collaborator
			farms[collaborator].Size += farms[v].Size
		}

		sig, err := signData(farms[v].Address, []byte(farms[staker].Address))
		if err != nil {
			if err.Error() == errors.New(fmt.Sprintf("Internal error: %v", facWallet.ErrNoSuchAddress)).Error() {
				return errors.New(fmt.Sprintf("[ERROR] invalid address, address : %v", farms[v].Address))
			}
			return errors.New(fmt.Sprintf("[ERROR] unable to sign data, err : %v", err))
		}

		farms[v].Content = farms[staker].Address
		farms[v].Sig, farms[v].PubKey = sig.Signature, sig.PubKey
	}

	return nil
}

func calculateReward(farms map[string]*farm, farmer string, divisor, cropSize int) (float64, error) {
	content := farms[farmer].Content
	if err := verifySig([]byte(content), farms[farmer].PubKey, farms[farmer].Sig); err != nil {
		return 0, err
	}
	return (float64(reward) / float64(divisor)) * float64(cropSize), nil
}

func signData(signer string, data []byte) (*fac.Signature, error) {
	sig, err := fac.SignData(signer, data)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func verifySig(data, pubKey, sig []byte) error {
	return primitives.VerifySignature(data, pubKey, sig)
}
