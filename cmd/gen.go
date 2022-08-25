package cmd

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gitlab.rete.farm/sistemi/inca/pki"
	"gitlab.rete.farm/sistemi/inca/util"
)

var cmdGen = &cobra.Command{
	Use:   "gen",
	Short: "Generate CA certificate",
	Run: func(cmd *cobra.Command, args []string) {
		output, err := cmd.Flags().GetString("output")
		if err != nil {
			log.Fatal().Err(err).Msg("output flag is mandatory")
		}
		if output == "-" {
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		}

		names, err := cmd.Flags().GetStringArray("name")
		if err != nil {
			log.Fatal().Err(err).Msg("at least a name gotta be given")
		}

		reqOptions := make(map[string]any)
		reqOptions["cn"] = names[0]
		if len(names) > 1 {
			reqOptions["alt"] = names[1:]
		}

		if duration, err := cmd.Flags().GetDuration("duration"); err == nil {
			reqOptions["duration"] = duration
		}

		if ca, err := cmd.Flags().GetBool("ca"); err == nil {
			reqOptions["ca"] = ca
		}

		for _, reqOptionKey := range []string{
			"algo", "organization", "country", "province", "locality", "street-address", "postal-code",
		} {
			if reqOptionValue, err := cmd.Flags().GetString(reqOptionKey); err == nil {
				reqOptions[strings.ReplaceAll(reqOptionKey, "-", "_")] = reqOptionValue
			}
		}

		req := pki.NewRequest(reqOptions)
		log.Info().Str("name", req.CN).
			Strs("dns", req.DNSNames).
			Strs("ip", req.IPAddresses).
			Dur("duration", req.Duration).
			Str("algo", string(req.Algo)).
			Msg("generating certificate")

		crt, key, err := pki.New(req)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to generate certificate")
		}
		log.Info().Str("certificate", crt.SerialNumber.String()).Msg("certificate generated")
		log.Info().Str("public", fmt.Sprintf("%v", key.Public())).Msg("key generated")

		log.Info().Msg("wrapping certificate")
		crtBytes, keyBytes, err := pki.Wrap(crt, key, crt, key)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to wrap certificate")
		}

		if output != "-" {

			log.Info().Msg("exporting certificate")
			if err := pki.Export(crtBytes, filepath.Join(output, "crt.pem")); err != nil {
				log.Fatal().Err(err).Msg("unable to export certificate")
			}

			log.Info().Msg("exporting key")
			if err := pki.Export(keyBytes, filepath.Join(output, "key.pem")); err != nil {
				log.Fatal().Err(err).Msg("unable to export key")
			}

			log.Info().Str("output", output).Msg("certificate created")
			return

		} else {

			var (
				crtBuffer = pki.ExportBytes(crtBytes)
				keyBuffer = pki.ExportBytes(keyBytes)
			)

			encode, err := cmd.Flags().GetString("encode")
			if err != nil {
				log.Fatal().Err(err).Msg("unable to read encode flag")
			}

			switch encode {
			case "zip":
				out := new(bytes.Buffer)
				zip := zip.NewWriter(out)
				for key, value := range map[string][]byte{
					"crt.pem": crtBuffer,
					"key.pem": keyBuffer,
				} {
					file, err := zip.Create(key)
					if err != nil {
						log.Fatal().Err(err).Msg("unable to create ZIP archive entry")
					}

					if _, err := file.Write(value); err != nil {
						log.Fatal().Err(err).Msg("unable to add content to ZIP archive entry")
					}
				}

				err := zip.Close()
				if err != nil {
					log.Fatal().Err(err).Msg("unable to close ZIP archive")
				}

				cmd.Print(out.String())
			case "json":
				jsonBundle := &struct {
					Crt string `json:"crt"`
					Key string `json:"key"`
				}{string(crtBuffer), string(keyBuffer)}
				if json, err := json.Marshal(jsonBundle); err != nil {
					log.Fatal().Err(err).Msg("unable to JSON-encode archive")
				} else {
					cmd.Print(string(json))
				}
			default: // raw
				cmd.Printf("%s%s", string(crtBuffer), string(keyBuffer))
			}
			return

		}
	},
}

func init() {
	cmdRoot.AddCommand(cmdGen)
	cmdGen.Flags().StringArrayP("name", "n", []string{}, "Certificate names")
	cmdGen.Flags().StringP("output", "o", util.ErrWrap("./")(os.Getwd()), "Output path (\"-\" for stdout)")
	cmdGen.Flags().StringP("encode", "e", "raw", "Encode returned payload: zip, json (only for stdout generation)")
	cmdGen.Flags().StringP("algo", "a", pki.DefaultCrtAlgo, "Private key algorithm")
	cmdGen.Flags().String("organization", "", "Certificate Organization")
	cmdGen.Flags().String("country", "", "Certificate Country")
	cmdGen.Flags().String("province", "", "Certificate Province")
	cmdGen.Flags().String("locality", "", "Certificate Locality")
	cmdGen.Flags().String("street-address", "", "Certificate StreetAddress")
	cmdGen.Flags().String("postal-code", "", "Certificate PostalCode")
	cmdGen.Flags().Duration("duration", pki.DefaultCrtDuration, "Certificate Duration")
	cmdGen.Flags().Bool("ca", false, "CA-enabled certificate")
	if err := cmdGen.MarkFlagRequired("name"); err != nil {
		log.Fatal().Err(err).Msg("unable to mark name flag as required")
	}
}
