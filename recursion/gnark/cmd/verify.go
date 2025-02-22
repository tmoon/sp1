package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/spf13/cobra"
	"github.com/succinctlabs/sp1-recursion-gnark/sp1"
	"github.com/succinctlabs/sp1-recursion-gnark/sp1/babybear"
)

var verifyCmdDataDir string
var verifyCmdProof string
var verifyCmdVkeyHash string
var verifyCmdCommitedValuesDigest string

func init() {
	verifyCmd.Flags().StringVar(&verifyCmdDataDir, "data", "", "")
	verifyCmd.Flags().StringVar(&verifyCmdProof, "proof", "", "")
	verifyCmd.Flags().StringVar(&verifyCmdVkeyHash, "vkey-hash", "", "")
	verifyCmd.Flags().StringVar(&verifyCmdCommitedValuesDigest, "commited-values-digest", "", "")
}

var verifyCmd = &cobra.Command{
	Use: "verify",
	Run: func(cmd *cobra.Command, args []string) {
		// Sanity check the required arguments have been provided.
		if verifyCmdDataDir == "" {
			panic("--data is required")
		}

		// Decode the proof.
		proofDecodedBytes, err := hex.DecodeString(verifyCmdProof)
		if err != nil {
			panic(err)
		}
		proof := groth16.NewProof(ecc.BN254)
		if _, err := proof.ReadFrom(bytes.NewReader(proofDecodedBytes)); err != nil {
			panic(err)
		}

		// Read the verifier key.
		vkFile, err := os.Open(verifyCmdDataDir + "/" + VK_PATH)
		if err != nil {
			panic(err)
		}
		vk := groth16.NewVerifyingKey(ecc.BN254)
		vk.ReadFrom(vkFile)

		// Compute the public witness.
		circuit := sp1.Circuit{
			Vars:                 []frontend.Variable{},
			Felts:                []babybear.Variable{},
			Exts:                 []babybear.ExtensionVariable{},
			VkeyHash:             verifyCmdVkeyHash,
			CommitedValuesDigest: verifyCmdCommitedValuesDigest,
		}
		witness, err := frontend.NewWitness(&circuit, ecc.BN254.ScalarField())
		if err != nil {
			panic(err)
		}
		publicWitness, err := witness.Public()
		if err != nil {
			panic(err)
		}

		// Verify proof.
		err = groth16.Verify(proof, vk, publicWitness, backend.WithVerifierHashToFieldFunction(sha256.New()))
		if err != nil {
			panic(err)
		}

	},
}
