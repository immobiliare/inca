## [1.8.1](https://github.com/immobiliare/inca/compare/1.8.0...1.8.1) (2025-07-07)


### Bug Fixes

* **oci:** create missing directory .well-known for webroot ([5eee4c0](https://github.com/immobiliare/inca/commit/5eee4c08075116d60fdf9efc30928f352a0663fd))
* **oci:** create missing subdirectory for acme-challenge in webroot ([4e63042](https://github.com/immobiliare/inca/commit/4e63042d439a9112a53afc0f58dc8e1f46ab68a8))
* **oci:** set permissions for webroot directory ([133596f](https://github.com/immobiliare/inca/commit/133596fee1fc6e54bdc63a0e220bbf42abfdf6c4))

## [1.8.0](https://github.com/immobiliare/inca/compare/1.7.2...1.8.0) (2025-07-02)


### Features

* **server:** add PFX download functionality ([96ef5dd](https://github.com/immobiliare/inca/commit/96ef5dd864cbee2d741dc4465e9ba701ee37c019))
* **server:** detect leaf and chain certs before encoding PFX ([84b6cad](https://github.com/immobiliare/inca/commit/84b6cad444aeb911b603ae44b333bb126341fc6d))
* **util:** add GenerateRandomString functionality ([f49b6b6](https://github.com/immobiliare/inca/commit/f49b6b699394ee6bb38bb136979fed6460aade90))


### Bug Fixes

* **oci:** update syft output path to an unprivileged directory ([4063d47](https://github.com/immobiliare/inca/commit/4063d47a470247e218d17701905628df19f8e6f0))
* **server:** add a nil check to handle invalid PEM content ([4fc24cd](https://github.com/immobiliare/inca/commit/4fc24cd8106dde0bfdf43993027124bca81ad72b))
* **server:** add a nil check to handle invalid PEM content ([89e3ab5](https://github.com/immobiliare/inca/commit/89e3ab5ad58184b02f7648e9fd8bfac29352bf9a))
* **server:** enhance private key parsing to support multiple formats ([0481710](https://github.com/immobiliare/inca/commit/04817108e99a9cadff44ced642d28bab06056f9f))
* **server:** reduce password length for PFX generation from 256 to 30 characters ([58582f1](https://github.com/immobiliare/inca/commit/58582f18c0359b349bb785768f8ea8f53dfdd8c9))

## [1.7.2](https://github.com/immobiliare/inca/compare/1.7.1...1.7.2) (2025-05-23)


### Bug Fixes

* **deps:** bump github.com/gofiber/fiber/v2 from 2.52.6 to 2.52.7 ([3a80c20](https://github.com/immobiliare/inca/commit/3a80c20ed393cb51744fb8bdf13f6eaab857e970))

## [1.7.1](https://github.com/immobiliare/inca/compare/1.7.0...1.7.1) (2025-05-06)


### Bug Fixes

* **cmd:** update command output method to SetOut for consistency ([82e9618](https://github.com/immobiliare/inca/commit/82e96181510ef1b59297d5ade3c8975917387034))
* **oci:** add missing provenance and sbom options to the docker build-push action ([b30e8ab](https://github.com/immobiliare/inca/commit/b30e8ab93312ecad7368fdccd488efaef3b29ca7))
* **oci:** update dockerfile to use chainguard base images ([b887b40](https://github.com/immobiliare/inca/commit/b887b408cac9e8b06aa61f9f0e1a50f64656c014))
* **pki:** ensure every error is handled in export function ([cbcd57b](https://github.com/immobiliare/inca/commit/cbcd57b5990243ffe033a0a88bcdc2bc8e9aea3e))
* **provider:** ensure every error is handled in ca function ([aac0854](https://github.com/immobiliare/inca/commit/aac0854cf436f15bb5b1f16727f728f0b379ba5f))
* **server:** replace unmaintained gopkg.in/yaml.v3 with github.com/goccy/go-yaml ([6463248](https://github.com/immobiliare/inca/commit/6463248024162e9286853ce3d5bd137cdd8b0397))
* **storage:** ensure every error is handled in find function ([63ef369](https://github.com/immobiliare/inca/commit/63ef3699cd20d876a1d47212aace6efdff53312a))
* **storage:** ensure every error is handled in get function ([dd3b299](https://github.com/immobiliare/inca/commit/dd3b29993754113039ba853f9cc85a8dbdf1bbff))
* **test:** handle error when writing response in httpRequestHandler ([1c2e662](https://github.com/immobiliare/inca/commit/1c2e66268af4e3e90725023e099ae48df55bf8e6))

## [1.7.0](https://github.com/immobiliare/inca/compare/1.6.0...1.7.0) (2025-01-07)


### Features

* **storage:** add postgresql storage implementation ([c26a105](https://github.com/immobiliare/inca/commit/c26a105d71aff801079acff5a28331bf466b9d04))


### Bug Fixes

* **http-01:** switch from relative to absolute paths ([ef43c90](https://github.com/immobiliare/inca/commit/ef43c90c1404a6b3e442fa60242f66b4a69663bc))
* **http-01:** the provider is creating the directory structure ([8362e1b](https://github.com/immobiliare/inca/commit/8362e1be287385c87ee6399a1333ee4b037183b1))

## [1.6.0](https://github.com/immobiliare/inca/compare/1.5.1...1.6.0) (2024-10-08)


### Features

* add mTLS support and ensure backward compatibility ([87fe565](https://github.com/immobiliare/inca/commit/87fe565d382889fc78887a6b8b2ef6543485722a))


### Bug Fixes

* **provider:** patch CA URL handling in LetsEncrypt Tune function ([d8f697c](https://github.com/immobiliare/inca/commit/d8f697c7700f26a705553d585b42456c02d995a2))
* **provider:** patch type assertion in LetsEncrypt.Tune method ([f7988aa](https://github.com/immobiliare/inca/commit/f7988aa31f0db12283cdc4f798b6a9c89f41031c))

## [1.5.1](https://github.com/immobiliare/inca/compare/1.5.0...1.5.1) (2024-07-02)


### Bug Fixes

* **security:** patch for a session middleware injection vulnerability in gofiber (CVE-2024-38513) ([1aaa0ce](https://github.com/immobiliare/inca/commit/1aaa0ceb38a3e8b3f81f7179fd065e18b775766a))

## [1.5.0](https://github.com/immobiliare/inca/compare/1.4.1...1.5.0) (2024-06-12)


### Features

* **docs:** introduce a documentation for the configuration yaml ([b371587](https://github.com/immobiliare/inca/commit/b371587b1f001bbbb783fb23cc415e7c2d90db61))


### Bug Fixes

* **deps:** manually bump the baseline - security patches ([95a05ca](https://github.com/immobiliare/inca/commit/95a05ca020ab164750cd572aa16369efcb727414))

## [1.4.1](https://github.com/immobiliare/inca/compare/1.4.0...1.4.1) (2024-03-08)


### Bug Fixes

* **storage:** missing validation for s3 bucket names ([b5a0dcc](https://github.com/immobiliare/inca/commit/b5a0dcc410be5e3574e790608d43a87589019641))

## [1.4.0](https://github.com/immobiliare/inca/compare/1.3.0...1.4.0) (2024-02-07)


### Features

* **ci:** introduce a conventional commits linter ([e851834](https://github.com/immobiliare/inca/commit/e8518348701ac4d2e16cb4ab4bcf33f02027833e))


### Bug Fixes

* **provider:** trim whitespaces in SAN fields [letsencrypt] ([11575d9](https://github.com/immobiliare/inca/commit/11575d94dea1364e2591c71bf35dfc9bec1eba0a))

## [1.3.0](https://github.com/immobiliare/inca/compare/1.2.0...1.3.0) (2024-01-08)


### Features

* **provider:** support for the HTTP-01 challenge ([#29](https://github.com/immobiliare/inca/issues/29)) ([f89dd6b](https://github.com/immobiliare/inca/commit/f89dd6bedbb1337431749de79d1d5919370caf6a))


### Bug Fixes

* **test:** switch the path for unit testing to an absolute one ([65492c5](https://github.com/immobiliare/inca/commit/65492c5891bb22502f982a759fbd9c8d45e67be1))

## [1.2.0](https://github.com/immobiliare/inca/compare/1.1.1...1.2.0) (2024-01-05)


### Features

* **ci:** generate release from github action and update changelog ([e19e433](https://github.com/immobiliare/inca/commit/e19e4331c436bf22afb8e4d113cd9970b4b8039e))
* **ci:** introduce an automatic dependency managament based on dependabot ([ae4d919](https://github.com/immobiliare/inca/commit/ae4d919c784df01329a0e9faab8f3e8ff92ac353))


### Bug Fixes

* **ci:** action release now it's ready to run ([4e1bd37](https://github.com/immobiliare/inca/commit/4e1bd37de9e0a896dacca594fbdbefc04f528609))

## [1.1.1](https://github.com/immobiliare/inca/compare/1.1.0...1.1.1) (2023-12-19)

### Bug Fixes

* **ci:** remove platform to test publish ([bda6273](https://github.com/immobiliare/inca/commit/bda62731acc3cdd4cc1d1fb4f315e15ccd7c9433))


## [1.1.0](https://github.com/immobiliare/inca/compare/1.0.8...1.1.0) (2023-12-19)

### Features

* General Project Maintenance: dependency updates and test fixes
* **test:** remove the expired embedded x509 keypair ([5fe8366](https://github.com/immobiliare/inca/commit/5fe836656d32fcac389e708358c89f48ada85eec))

### Bug Fixes

* **ci:** add platforms also in tests ([c62354c](https://github.com/immobiliare/inca/commit/c62354cb166e23f8e5cf9a59833f26ccac91e3b0))


## [1.0.8](https://github.com/immobiliare/inca/compare/1.0.7...1.0.8) (2023-12-19)

### Bug Fixes
* **ci:** try new publish release ([a6073f9](https://github.com/immobiliare/inca/commit/a6073f9212168d0d6cee319678e3dd1effc8bc14))


## [1.0.7](https://github.com/immobiliare/inca/compare/1.0.6...1.0.7) (2023-12-19)

### Bug Fixes

* **ci:** remove new tag env ([22f84dc](https://github.com/immobiliare/inca/commit/22f84dcf96ebb44974bb07e8645a5468b90cf037))


## [1.0.6](https://github.com/immobiliare/inca/compare/1.0.5...1.0.6) (2023-12-19)


### Features

* **server:** add keypair download button ([be81761](https://github.com/immobiliare/inca/commit/be817617d34bd2c00df54b0682b6129e815e42aa))


### Bug Fixes

* ensure SANs is filled with Subject CN too ([63012c1](https://github.com/immobiliare/inca/commit/63012c17f149c9514b4c8dd0e767924d9f77590e))

## [1.0.4](https://github.com/immobiliare/inca/compare/1.0.3...1.0.4) (2022-10-11)

## [1.0.3](https://github.com/immobiliare/inca/compare/1.0.2...1.0.3) (2022-10-11)


### Bug Fixes

* **pki:** add support for parsing RSA-encoded keys ([56e0ec8](https://github.com/immobiliare/inca/commit/56e0ec8a4b4eacbcf5fbfe4af4a55aca625a4a9e))

## [1.0.2](https://github.com/immobiliare/inca/compare/1.0.1...1.0.2) (2022-10-11)

## [1.0.1](https://github.com/immobiliare/inca/compare/1.0.0...1.0.1) (2022-10-10)


### Features

* **server:** enable certificate autorenew on get ([96dcd46](https://github.com/immobiliare/inca/commit/96dcd4621e0154265e3615d833ebfe2b6b0ad463))


### Bug Fixes

* **web:** do not include session flows in non-acl mode ([e18e572](https://github.com/immobiliare/inca/commit/e18e572d532c4c419e1b2a1805aed0d84ad893bb))

## [1.0.0](https://github.com/immobiliare/inca/compare/0.1.0...1.0.0) (2022-10-07)


### Features

* **acl:** authorization-based API access plus token based web session ([3fb84c3](https://github.com/immobiliare/inca/commit/3fb84c31afd15114f285aaee74c922b86929ecbd))
* **acl:** implement target-based authorization scheme ([9a69969](https://github.com/immobiliare/inca/commit/9a6996944c731180d27538ac005512cec13431e5))
* **acl:** unprotect on empty ACLs ([9f61dd8](https://github.com/immobiliare/inca/commit/9f61dd8b3b115cd6f798f22ff5fc22194b102ec4))
* add certificate revocation at the provider level (j#IS-3039) ([0423171](https://github.com/immobiliare/inca/commit/0423171bc31155fbff81c29e9bee5c9bed8c29df)), closes [j#IS-3039](https://github.com/immobiliare/j/issues/IS-3039)
* add not-logged health endpoint ([bc45cb5](https://github.com/immobiliare/inca/commit/bc45cb5d70ba56c7619e86dd518d9782a4f45b02))
* add support for querying same-type CAs ([76216dc](https://github.com/immobiliare/inca/commit/76216dc4c2229c500674fe127b876cf952bafbe5))
* **docker:** switch to entrypoint from cmd ([757bafe](https://github.com/immobiliare/inca/commit/757bafeefb9a4058c229daa0967c0e042380e31f))
* **gen:** add support for stdout/compress CA generation ([9c11195](https://github.com/immobiliare/inca/commit/9c11195ff737ef1794ac7809db1f1a3faeb863c5))
* **gen:** replace `compress` flag with `encode` and add json support ([6f41950](https://github.com/immobiliare/inca/commit/6f419500167c7a6782fbc2796cf01f081d879fac))
* **letsencrypt:** first dump ([af36645](https://github.com/immobiliare/inca/commit/af36645313b256721e4206fb6a713057588b48f4))
* **local:** add support for certificates alt names ([deb9677](https://github.com/immobiliare/inca/commit/deb9677418439707efcbaada2626a4bf7f9b5fde))
* **local:** add support for custom key algorithm ([ed70157](https://github.com/immobiliare/inca/commit/ed7015761a8668aea9a116cf097b4f6aa36921c4))
* merge alt names into existing certificates (j#IS-2865) ([a35b89d](https://github.com/immobiliare/inca/commit/a35b89d567631fc681cb3114d5bda80ea4aa757e)), closes [j#IS-2865](https://github.com/immobiliare/j/issues/IS-2865)
* **pki:** add support for parsing EC-encoded keys ([66d98f1](https://github.com/immobiliare/inca/commit/66d98f154c8fcb18228c8b4bca962adc9cb15e0e))
* **pki:** switch to ECDSA with SHA-256 ([5678544](https://github.com/immobiliare/inca/commit/5678544ff005a4d7889b69e837d069af8dee6ca3))
* rework config parsing to make it stateful (j#IS-2874) ([be04f9f](https://github.com/immobiliare/inca/commit/be04f9ffae73efc50038b05b825f4c7289c13135)), closes [j#IS-2874](https://github.com/immobiliare/j/issues/IS-2874)
* **sentry:** first dump ([0af58a6](https://github.com/immobiliare/inca/commit/0af58a6fad9644b7f0f5c8f2d1cf16978346c4d0))
* **server:** add certificate show endpoint ([959b4a3](https://github.com/immobiliare/inca/commit/959b4a3a105a8319ecbf3b66be93c90e8cd7ec08))
* **server:** add compression middleware ([c2d7c2d](https://github.com/immobiliare/inca/commit/c2d7c2deac46b701657f97d05c3094e4c54e75d7))
* **server:** expose endpoint for certificates enumeration (j#IS-2824) ([c43b530](https://github.com/immobiliare/inca/commit/c43b530a1d3cfd6fe1f7a507b5960f6a0e698a66)), closes [j#IS-2824](https://github.com/immobiliare/j/issues/IS-2824)
* **server:** support JSON encoding as an option ([02ad193](https://github.com/immobiliare/inca/commit/02ad1934d50ba84d2fbd2661cd88e241e167fd91))
* **storage:** add support for S3 ([5f1bcd9](https://github.com/immobiliare/inca/commit/5f1bcd9579a1af2cd625d03e56f4e2d4ee051c1e))
* **storage:** treat assets as unicum ([6c70f2e](https://github.com/immobiliare/inca/commit/6c70f2e2331f39316a1b71c2d8458f834febb830))
* **webgui:** first dump ([ee3df98](https://github.com/immobiliare/inca/commit/ee3df9805a4995e3ed02f91fddb9436f8773b373))


### Bug Fixes

* **ca:** properly bundle CA certificate ([e839bfd](https://github.com/immobiliare/inca/commit/e839bfd97ff4908f8a621cfcb2f4280801f171f5))
* debrand hardcoded certificates defaults (j#IS-2853) ([cc88f5a](https://github.com/immobiliare/inca/commit/cc88f5a2bfca90038b2e845e99e942f1db1c0fe9)), closes [j#IS-2853](https://github.com/immobiliare/j/issues/IS-2853)
* **docker:** copy static assets ([d7aaaf9](https://github.com/immobiliare/inca/commit/d7aaaf9cd299f3620e2cffa0d7f865fcfd492ad4))
* **gen:** use first name as CN ([17f0fae](https://github.com/immobiliare/inca/commit/17f0fae21caf3d9520d53e5cd1e023a17dbf98b1))
* **gen:** use unix-like flag names ([c2560ac](https://github.com/immobiliare/inca/commit/c2560acb8ab8770884803f700ddc4945c89377e8))
* **letsencrypt:** proxy CA certificate retrieval ([cbfe75e](https://github.com/immobiliare/inca/commit/cbfe75e7684370a6e34ef28c6e134d0e81be138c))
* **local:** check suffix match on CN too ([0252d14](https://github.com/immobiliare/inca/commit/0252d1469910aec5565b47c9ba80b5dc6545d0a0))
* **pki:** properly check CN as domain name ([26d7708](https://github.com/immobiliare/inca/commit/26d770812d083332fcab0daa289e1b35c01a27be))
* **pki:** reduce default crt duration to 397 days (j#IS-2964) ([36ff40b](https://github.com/immobiliare/inca/commit/36ff40b000931e07a619c8b634ddb777ea0d80a8)), closes [j#IS-2964](https://github.com/immobiliare/j/issues/IS-2964)
* **pki:** support 4+ level domain names ([1eb80be](https://github.com/immobiliare/inca/commit/1eb80bec809fc5b12c017f34161d670b3540acff))
* **pki:** use default algo constant in request generator ([996eb4c](https://github.com/immobiliare/inca/commit/996eb4cf0a248f3c9771f6bbcf77fc0c2ecdd98a))
* **pki:** use url-compliant variable names for crt parameters ([e2a763f](https://github.com/immobiliare/inca/commit/e2a763f4ba4f9cff17494bb9c1e46ca75791e76a))
* **s3:** do not wait for bucket ([3d54fee](https://github.com/immobiliare/inca/commit/3d54fee3440e69b0c45a0bb7faf9f89702afe657))
* **s3:** rework S3 logic ([9becb58](https://github.com/immobiliare/inca/commit/9becb58996d13ada02791e992aa2201285b6d869))
* **server:** API already ACL-protected at the handler-level ([038cdf6](https://github.com/immobiliare/inca/commit/038cdf63da7154c53f7e2ccccaf1ab611beb19b1))
* **server:** prefer using c.JSON where returning a JSON-encoded responses (j#IS-2863) ([ab7fe23](https://github.com/immobiliare/inca/commit/ab7fe230d7bb589de59b1d05442ed1350fc2ebf7)), closes [j#IS-2863](https://github.com/immobiliare/j/issues/IS-2863)
* **server:** use InternalServerError instead of BadRequest when suitable (j#IS-2864) ([ebb3ba1](https://github.com/immobiliare/inca/commit/ebb3ba1ec5ac08da1ad36a72f70c49e72e376b74)), closes [j#IS-2864](https://github.com/immobiliare/j/issues/IS-2864)
* **zstdlogger:** restore correct default output fd ([f30ab18](https://github.com/immobiliare/inca/commit/f30ab183bdb7899a0b96220a4454d2c04fc84d90))

## [0.1.0](https://github.com/immobiliare/inca/compare/58340280298126d2434f36bcca07aa13c8802768...0.1.0) (2022-07-15)


### Features

* add certificates storing/returning to/from storage ([62b2ecd](https://github.com/immobiliare/inca/commit/62b2ecdd2acada73dc03ac3897cefe50f2379a53))
* add support for certificates removal ([c6fb40c](https://github.com/immobiliare/inca/commit/c6fb40c3cf5be12a37259530a97c49ae321b1834))
* add support for docker ([0c91463](https://github.com/immobiliare/inca/commit/0c91463bce969b70bea02d8d38e5e679826c6de5))
* expose configurable providers and add certficate parsing ([d47d358](https://github.com/immobiliare/inca/commit/d47d358210729cb1158d417fa9855cdb5c4227b1))
* expose endpoint for CA certificate ([624ff6d](https://github.com/immobiliare/inca/commit/624ff6d2b550880ec1bb611fd4aa8d1823f1979d))
* get a basic functional crt/key pair generation logic ([5834028](https://github.com/immobiliare/inca/commit/58340280298126d2434f36bcca07aa13c8802768))
* **server:** add support for certificates creation plus ECDSA algo ([95a4fd4](https://github.com/immobiliare/inca/commit/95a4fd4b07d537fda79d69a32d8b6b84c6f2ff8e))


### Bug Fixes

* **docker:** run server only ([c5388b2](https://github.com/immobiliare/inca/commit/c5388b200616023da06c75e74166409e06277e5b))
* **gen:** use defaults-filled certificate request ([1873404](https://github.com/immobiliare/inca/commit/1873404e218d8a7643812e8bd66f4d2d661a4d58))
* **pki:** export certificates in the given path ([c7df8c2](https://github.com/immobiliare/inca/commit/c7df8c251ebce75f5d09d1cc46649616e911e8f4))
