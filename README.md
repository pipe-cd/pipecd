[![Build](https://github.com/pipe-cd/pipecd/actions/workflows/build.yaml/badge.svg)](https://github.com/pipe-cd/pipecd/actions/workflows/build.yaml)
[![Test](https://github.com/pipe-cd/pipecd/actions/workflows/test.yaml/badge.svg)](https://github.com/pipe-cd/pipecd/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/pipe-cd/pipecd)](https://goreportcard.com/report/github.com/pipe-cd/pipecd)
[![codecov](https://codecov.io/gh/pipe-cd/pipecd/branch/master/graph/badge.svg)](https://codecov.io/gh/pipe-cd/pipecd)
[![Release](https://img.shields.io/github/v/release/pipe-cd/pipecd?label=Release)](https://github.com/pipe-cd/pipecd/releases/latest)
[![License](https://img.shields.io/badge/License-Apache_2.0-blue.svg)](/LICENSE)
[![Documentation](https://img.shields.io/badge/Documentation-pipecd-informational.svg)](https://pipecd.dev/docs/)
[![Slack](https://img.shields.io/badge/Slack-%23pipecd-informational.svg)](https://app.slack.com/client/T08PSQ7BQ/C01B27F9T0X)
[![Twitter](https://img.shields.io/twitter/url/https/twitter.com/pipecd_dev.svg?style=social&label=Follow%20%40pipecd_dev)](https://twitter.com/pipecd_dev)

<p align="center">
  <img src="https://github.com/pipe-cd/pipecd/blob/master/docs/static/images/logo.png" width="180"/>
</p>

<p align="center">
  A GitOps style continuous delivery platform that provides consistent deployment and operations experience for any applications
  <br/>
  <a href="https://pipecd.dev"><strong>Explore PipeCD docs »</strong></a>
  <a href="https://play.pipecd.dev?project=play"><strong>Play with live demo »</strong></a>
</p>

#

![](https://github.com/pipe-cd/pipecd/blob/master/docs/static/images/rolled-back-deployment.png)

### Overview

PipeCD provides __a unified continuous delivery solution for multiple application kinds on multi-cloud__ that empowers engineers to deploy faster with more confidence, a GitOps tool that enables doing deployment operations by pull request on Git.

![](https://github.com/pipe-cd/pipecd/blob/master/docs/static/images/pipecd-explanation.png)

#

### Why PipeCD?

- Simple, unified and easy to use but powerful pipeline definition to construct your deployment
- Same deployment interface to deploy applications of any platform, including Kubernetes, Terraform, GCP Cloud Run, AWS Lambda, AWS ECS
- No CRD or applications' manifest changes are required; Only need a pipeline definition along with your application manifests
- No deployment credentials are exposed or required outside the application cluster
- Built-in deployment analysis as part of the deployment pipeline to measure impact based on metrics, logs, emitted requests
- Easy to interact with any CI; The CI tests and builds artifacts, PipeCD takes the rest
- Insights show metrics like lead time, deployment frequency, MTTR and change failure rate to measure delivery performance
- Designed to manage thousands of cross-platform applications in multi-cloud for company scale but also work well for small projects

And many more, please explore in [docs](https://pipecd.dev/docs).

#

### Quickstart

The [quickstart guide](https://pipecd.dev/docs/quickstart/) shows how to set up PipeCD components and deploy a hello-world application with PipeCD.

#

### Installation

The [installation guide](https://pipecd.dev/docs/installation/) explains and helps set up PipeCD for your real-life production environment.

#

### Community and development

- Check out the [PipeCD website](https://pipecd.dev) for the complete documentation and helpful links.
- Join our [Slack workspace](https://cloud-native.slack.com/archives/C01B27F9T0X) to get help and to discuss features.
- Tweet [@pipecd_dev](https://twitter.com/pipecd_dev) on Twitter/X.
- Create Github [Issues](https://github.com/pipe-cd/pipecd/issues) or [Discussions](https://github.com/pipe-cd/pipecd/discussions/) to report bugs or request features.
- Join our PipeCD Development and Community Meeting where we share the latest project news, demos, answer questions, and triage issues.
  - 🗓️ Bi-weekly (twice a month)
  - ⏰ 14:30 Hanoi | 16:30 Tokyo | 15:30 Shanghai | 13:00 India | 8:30 CET | 23:30 PST
  - 🗒️ [Meeting minutes](https://bit.ly/pipecd-mtg-notes) (zoom link and schedule included)

Participation in PipeCD project is governed by the CNCF [Code of Conduct](CODE_OF_CONDUCT.md).

#

### Contributing

We'd love you to join us! Please see the [Contributing docs](CONTRIBUTING.md) to get started.

#

### Adopters

You can find a list of publicly recognized users of the PipeCD in the [ADOPTERS.md](ADOPTERS.md) file. We strongly encourage all PipeCD users to add their names to this list, as we love to see the community's growing success!

#

### Thanks to the contributors of PipeCD!

<a href="https://github.com/nghialv"><img src="https://avatars.githubusercontent.com/u/1751755?v=4" title="nghialv" width="80" height="80"></a>
<a href="https://github.com/khanhtc1202"><img src="https://avatars.githubusercontent.com/u/32532742?v=4" title="khanhtc1202" width="80" height="80"></a>
<a href="https://github.com/nakabonne"><img src="https://avatars.githubusercontent.com/u/19730728?v=4" title="nakabonne" width="80" height="80"></a>
<a href="https://github.com/cakecatz"><img src="https://avatars.githubusercontent.com/u/6136383?v=4" title="cakecatz" width="80" height="80"></a>
<a href="https://github.com/knanao"><img src="https://avatars.githubusercontent.com/u/50069775?v=4" title="knanao" width="80" height="80"></a>
<a href="https://github.com/ono-max"><img src="https://avatars.githubusercontent.com/u/59436572?v=4" title="ono-max" width="80" height="80"></a>
<a href="https://github.com/kentakozuka"><img src="https://avatars.githubusercontent.com/u/16733673?v=4" title="kentakozuka" width="80" height="80"></a>
<a href="https://github.com/stormcat24"><img src="https://avatars.githubusercontent.com/u/919840?v=4" title="stormcat24" width="80" height="80"></a>
<a href="https://github.com/apps/dependabot"><img src="https://avatars.githubusercontent.com/in/29110?v=4" title="dependabot[bot]" width="80" height="80"></a>
<a href="https://github.com/Hosshii"><img src="https://avatars.githubusercontent.com/u/49914427?v=4" title="Hosshii" width="80" height="80"></a>
<a href="https://github.com/funera1"><img src="https://avatars.githubusercontent.com/u/60760935?v=4" title="funera1" width="80" height="80"></a>
<a href="https://github.com/ffjlabo"><img src="https://avatars.githubusercontent.com/u/40124947?v=4" title="ffjlabo" width="80" height="80"></a>
<a href="https://github.com/sanposhiho"><img src="https://avatars.githubusercontent.com/u/44139130?v=4" title="sanposhiho" width="80" height="80"></a>
<a href="https://github.com/Szarny"><img src="https://avatars.githubusercontent.com/u/26561120?v=4" title="Szarny" width="80" height="80"></a>
<a href="https://github.com/t-kikuc"><img src="https://avatars.githubusercontent.com/u/97105818?v=4" title="t-kikuc" width="80" height="80"></a>
<a href="https://github.com/sivchari"><img src="https://avatars.githubusercontent.com/u/55221074?v=4" title="sivchari" width="80" height="80"></a>
<a href="https://github.com/kevin-namba"><img src="https://avatars.githubusercontent.com/u/68955641?v=4" title="kevin-namba" width="80" height="80"></a>
<a href="https://github.com/hungran"><img src="https://avatars.githubusercontent.com/u/26101787?v=4" title="hungran" width="80" height="80"></a>
<a href="https://github.com/kurochan"><img src="https://avatars.githubusercontent.com/u/591247?v=4" title="kurochan" width="80" height="80"></a>
<a href="https://github.com/TaKO8Ki"><img src="https://avatars.githubusercontent.com/u/41065217?v=4" title="TaKO8Ki" width="80" height="80"></a>
<a href="https://github.com/chaspy"><img src="https://avatars.githubusercontent.com/u/10370988?v=4" title="chaspy" width="80" height="80"></a>
<a href="https://github.com/gkuga"><img src="https://avatars.githubusercontent.com/u/33643470?v=4" title="gkuga" width="80" height="80"></a>
<a href="https://github.com/karamaru-alpha"><img src="https://avatars.githubusercontent.com/u/38310693?v=4" title="karamaru-alpha" width="80" height="80"></a>
<a href="https://github.com/nnnkkk7"><img src="https://avatars.githubusercontent.com/u/68233204?v=4" title="nnnkkk7" width="80" height="80"></a>
<a href="https://github.com/TonkyH"><img src="https://avatars.githubusercontent.com/u/50762864?v=4" title="TonkyH" width="80" height="80"></a>
<a href="https://github.com/golemiso"><img src="https://avatars.githubusercontent.com/u/3282656?v=4" title="golemiso" width="80" height="80"></a>
<a href="https://github.com/ouchi2501"><img src="https://avatars.githubusercontent.com/u/11391317?v=4" title="ouchi2501" width="80" height="80"></a>
<a href="https://github.com/tnqv"><img src="https://avatars.githubusercontent.com/u/23372024?v=4" title="tnqv" width="80" height="80"></a>
<a href="https://github.com/seipan"><img src="https://avatars.githubusercontent.com/u/88176012?v=4" title="seipan" width="80" height="80"></a>
<a href="https://github.com/anhpnv"><img src="https://avatars.githubusercontent.com/u/40441000?v=4" title="anhpnv" width="80" height="80"></a>
<a href="https://github.com/arabian9ts"><img src="https://avatars.githubusercontent.com/u/24448137?v=4" title="arabian9ts" width="80" height="80"></a>
<a href="https://github.com/kanata2"><img src="https://avatars.githubusercontent.com/u/7460883?v=4" title="kanata2" width="80" height="80"></a>
<a href="https://github.com/na-ga"><img src="https://avatars.githubusercontent.com/u/537006?v=4" title="na-ga" width="80" height="80"></a>
<a href="https://github.com/TakumaKurosawa"><img src="https://avatars.githubusercontent.com/u/39955827?v=4" title="TakumaKurosawa" width="80" height="80"></a>
<a href="https://github.com/khanhtc3010"><img src="https://avatars.githubusercontent.com/u/9603918?v=4" title="khanhtc3010" width="80" height="80"></a>
<a href="https://github.com/tennashi"><img src="https://avatars.githubusercontent.com/u/10219626?v=4" title="tennashi" width="80" height="80"></a>
<a href="https://github.com/ShotaKitazawa"><img src="https://avatars.githubusercontent.com/u/19530785?v=4" title="ShotaKitazawa" width="80" height="80"></a>
<a href="https://github.com/gotyoooo"><img src="https://avatars.githubusercontent.com/u/6133219?v=4" title="gotyoooo" width="80" height="80"></a>
<a href="https://github.com/seitarof"><img src="https://avatars.githubusercontent.com/u/51070449?v=4" title="seitarof" width="80" height="80"></a>
<a href="https://github.com/minhquang4334"><img src="https://avatars.githubusercontent.com/u/22545967?v=4" title="minhquang4334" width="80" height="80"></a>
<a href="https://github.com/p0tr3c"><img src="https://avatars.githubusercontent.com/u/12850042?v=4" title="p0tr3c" width="80" height="80"></a>
<a href="https://github.com/Abirdcfly"><img src="https://avatars.githubusercontent.com/u/5100555?v=4" title="Abirdcfly" width="80" height="80"></a>
<a href="https://github.com/SakataAtsuki"><img src="https://avatars.githubusercontent.com/u/58636635?v=4" title="SakataAtsuki" width="80" height="80"></a>
<a href="https://github.com/butterv"><img src="https://avatars.githubusercontent.com/u/15773082?v=4" title="butterv" width="80" height="80"></a>
<a href="https://github.com/mura-s"><img src="https://avatars.githubusercontent.com/u/4702673?v=4" title="mura-s" width="80" height="80"></a>
<a href="https://github.com/Warashi"><img src="https://avatars.githubusercontent.com/u/3600530?v=4" title="Warashi" width="80" height="80"></a>
<a href="https://github.com/BIwashi"><img src="https://avatars.githubusercontent.com/u/49979368?v=4" title="BIwashi" width="80" height="80"></a>
<a href="https://github.com/moko-poi"><img src="https://avatars.githubusercontent.com/u/55864094?v=4" title="moko-poi" width="80" height="80"></a>
<a href="https://github.com/Linutux"><img src="https://avatars.githubusercontent.com/u/435352?v=4" title="Linutux" width="80" height="80"></a>
<a href="https://github.com/hosht"><img src="https://avatars.githubusercontent.com/u/3858627?v=4" title="hosht" width="80" height="80"></a>
<a href="https://github.com/ww24"><img src="https://avatars.githubusercontent.com/u/695166?v=4" title="ww24" width="80" height="80"></a>
<a href="https://github.com/tnir"><img src="https://avatars.githubusercontent.com/u/10229505?v=4" title="tnir" width="80" height="80"></a>
<a href="https://github.com/yoiki"><img src="https://avatars.githubusercontent.com/u/39365493?v=4" title="yoiki" width="80" height="80"></a>
<a href="https://github.com/JohnTitor"><img src="https://avatars.githubusercontent.com/u/25030997?v=4" title="JohnTitor" width="80" height="80"></a>
<a href="https://github.com/testwill"><img src="https://avatars.githubusercontent.com/u/8717479?v=4" title="testwill" width="80" height="80"></a>
<a href="https://github.com/mugioka"><img src="https://avatars.githubusercontent.com/u/62197019?v=4" title="mugioka" width="80" height="80"></a>
<a href="https://github.com/naonao2323"><img src="https://avatars.githubusercontent.com/u/74669884?v=4" title="naonao2323" width="80" height="80"></a>
<a href="https://github.com/tatsuya0429"><img src="https://avatars.githubusercontent.com/u/29541999?v=4" title="tatsuya0429" width="80" height="80"></a>
<a href="https://github.com/tokku5552"><img src="https://avatars.githubusercontent.com/u/69064290?v=4" title="tokku5552" width="80" height="80"></a>
<a href="https://github.com/fazledyn-or"><img src="https://avatars.githubusercontent.com/u/138655107?v=4" title="fazledyn-or" width="80" height="80"></a>
<a href="https://github.com/pheianox"><img src="https://avatars.githubusercontent.com/u/72671586?v=4" title="pheianox" width="80" height="80"></a>
<a href="https://github.com/ductnn"><img src="https://avatars.githubusercontent.com/u/22121217?v=4" title="ductnn" width="80" height="80"></a>
<a href="https://github.com/Forget-C"><img src="https://avatars.githubusercontent.com/u/25631506?v=4" title="Forget-C" width="80" height="80"></a>
<a href="https://github.com/hongchaodeng"><img src="https://avatars.githubusercontent.com/u/920884?v=4" title="hongchaodeng" width="80" height="80"></a>
<a href="https://github.com/hori-ryota"><img src="https://avatars.githubusercontent.com/u/2936501?v=4" title="hori-ryota" width="80" height="80"></a>
<a href="https://github.com/eltociear"><img src="https://avatars.githubusercontent.com/u/22633385?v=4" title="eltociear" width="80" height="80"></a>
<a href="https://github.com/sano307"><img src="https://avatars.githubusercontent.com/u/12808316?v=4" title="sano307" width="80" height="80"></a>
<a href="https://github.com/misukuro"><img src="https://avatars.githubusercontent.com/u/1040546?v=4" title="misukuro" width="80" height="80"></a>
<a href="https://github.com/masaaania"><img src="https://avatars.githubusercontent.com/u/2755429?v=4" title="masaaania" width="80" height="80"></a>
<a href="https://github.com/KeisukeYamashita"><img src="https://avatars.githubusercontent.com/u/23056537?v=4" title="KeisukeYamashita" width="80" height="80"></a>
<a href="https://github.com/Lennie"><img src="https://avatars.githubusercontent.com/u/330102?v=4" title="Lennie" width="80" height="80"></a>
<a href="https://github.com/gitbluf"><img src="https://avatars.githubusercontent.com/u/22802784?v=4" title="gitbluf" width="80" height="80"></a>
<a href="https://github.com/mi11km"><img src="https://avatars.githubusercontent.com/u/54844746?v=4" title="mi11km" width="80" height="80"></a>
<a href="https://github.com/jacobnguyenn"><img src="https://avatars.githubusercontent.com/u/60810674?v=4" title="jacobnguyenn" width="80" height="80"></a>
<a href="https://github.com/GotoRen"><img src="https://avatars.githubusercontent.com/u/63791288?v=4" title="GotoRen" width="80" height="80"></a>
<a href="https://github.com/RikiyaFujii"><img src="https://avatars.githubusercontent.com/u/23261497?v=4" title="RikiyaFujii" width="80" height="80"></a>

#

**We are a [Cloud Native Computing Foundation](https://cncf.io/) sandbox project.**

<img src="https://www.cncf.io/wp-content/uploads/2022/07/cncf-color-bg.svg" width=300 />

The Linux Foundation® (TLF) has registered trademarks and uses trademarks. For a list of TLF trademarks, see [Trademark Usage](https://www.linuxfoundation.org/trademark-usage/).
