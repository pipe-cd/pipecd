[![Build](https://github.com/pipe-cd/pipecd/actions/workflows/build.yaml/badge.svg)](https://github.com/pipe-cd/pipecd/actions/workflows/build.yaml)
[![Test](https://github.com/pipe-cd/pipecd/actions/workflows/test.yaml/badge.svg)](https://github.com/pipe-cd/pipecd/actions/workflows/test.yaml)
[![Release](https://img.shields.io/github/v/release/pipe-cd/pipecd?label=Release)](https://github.com/pipe-cd/pipecd/releases/latest)
[![Documentation](https://img.shields.io/badge/Documentation-pipecd-informational.svg)](https://pipecd.dev/docs/)
[![Slack](https://img.shields.io/badge/Slack-%23pipecd-informational.svg)](https://app.slack.com/client/T08PSQ7BQ/C01B27F9T0X)

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

![](https://github.com/pipe-cd/pipecd/blob/master/docs/static/images/deployment-details.png)

### Overview

PipeCD provides a unified continuous delivery solution for multiple application kinds on multi-cloud that empowers engineers to deploy faster with more confidence, a GitOps tool that enables doing deployment operations by pull request on Git.

![](https://github.com/pipe-cd/pipecd/blob/master/docs/static/images/pipecd-explanation.png)

**Visibility**
- Deployment pipeline UI shows clarify what is happening
- Separate logs viewer for each individual deployment
- Realtime visualization of application state
- Deployment notifications to slack, webhook endpoints
- Insights show metrics like lead time, deployment frequency, MTTR and change failure rate to measure delivery performance

**Automation**
- Automated deployment analysis to measure deployment impact based on metrics, logs, emitted requests
- Automatically roll back to the previous state as soon as analysis or a pipeline stage fails
- Automatically detect configuration drift to notify and render the changes
- Automatically trigger a new deployment when a defined event has occurred (e.g. container image pushed, helm chart published, etc)

**Safety and Security**
- Support single sign-on and role-based access control
- Credentials are not exposed outside the cluster and not saved in the control-plane
- Piped makes only outbound requests and can run inside a restricted network
- Built-in secrets management

**Multi-provider & Multi-Tenancy**
- Support multiple application kinds on multi-cloud including Kubernetes, Terraform, Cloud Run, AWS Lambda
- Support multiple analysis providers including Prometheus, Datadog, Stackdriver, and more
- Easy to operate multi-cluster, multi-tenancy by separating control-plane and piped

#

### License

Apache License 2.0, see [LICENSE](https://github.com/pipe-cd/pipecd/blob/master/LICENSE).

#

### Contributing

We'd love you to join us! Please see the [Contributor Guide](https://pipecd.dev/docs/contribution-guidelines/).

#

### Thanks to the contributors of PipeCD!

<a href="https://github.com/nghialv"><img src="https://avatars.githubusercontent.com/u/1751755?v=4" title="nghialv" width="80" height="80"></a>
<a href="https://github.com/khanhtc1202"><img src="https://avatars.githubusercontent.com/u/32532742?v=4" title="khanhtc1202" width="80" height="80"></a>
<a href="https://github.com/nakabonne"><img src="https://avatars.githubusercontent.com/u/19730728?v=4" title="nakabonne" width="80" height="80"></a>
<a href="https://github.com/cakecatz"><img src="https://avatars.githubusercontent.com/u/6136383?v=4" title="cakecatz" width="80" height="80"></a>
<a href="https://github.com/knanao"><img src="https://avatars.githubusercontent.com/u/50069775?v=4" title="knanao" width="80" height="80"></a>
<a href="https://github.com/ono-max"><img src="https://avatars.githubusercontent.com/u/59436572?v=4" title="ono-max" width="80" height="80"></a>
<a href="https://github.com/stormcat24"><img src="https://avatars.githubusercontent.com/u/919840?v=4" title="stormcat24" width="80" height="80"></a>
<a href="https://github.com/Hosshii"><img src="https://avatars.githubusercontent.com/u/49914427?v=4" title="Hosshii" width="80" height="80"></a>
<a href="https://github.com/sanposhiho"><img src="https://avatars.githubusercontent.com/u/44139130?v=4" title="sanposhiho" width="80" height="80"></a>
<a href="https://github.com/apps/dependabot"><img src="https://avatars.githubusercontent.com/in/29110?v=4" title="dependabot[bot]" width="80" height="80"></a>
<a href="https://github.com/Szarny"><img src="https://avatars.githubusercontent.com/u/26561120?v=4" title="Szarny" width="80" height="80"></a>
<a href="https://github.com/funera1"><img src="https://avatars.githubusercontent.com/u/60760935?v=4" title="funera1" width="80" height="80"></a>
<a href="https://github.com/TaKO8Ki"><img src="https://avatars.githubusercontent.com/u/41065217?v=4" title="TaKO8Ki" width="80" height="80"></a>
<a href="https://github.com/chaspy"><img src="https://avatars.githubusercontent.com/u/10370988?v=4" title="chaspy" width="80" height="80"></a>
<a href="https://github.com/ffjlabo"><img src="https://avatars.githubusercontent.com/u/40124947?v=4" title="ffjlabo" width="80" height="80"></a>
<a href="https://github.com/gkuga"><img src="https://avatars.githubusercontent.com/u/33643470?v=4" title="gkuga" width="80" height="80"></a>
<a href="https://github.com/kurochan"><img src="https://avatars.githubusercontent.com/u/591247?v=4" title="kurochan" width="80" height="80"></a>
<a href="https://github.com/TonkyH"><img src="https://avatars.githubusercontent.com/u/50762864?v=4" title="TonkyH" width="80" height="80"></a>
<a href="https://github.com/kevin55156"><img src="https://avatars.githubusercontent.com/u/68955641?v=4" title="kevin55156" width="80" height="80"></a>
<a href="https://github.com/tnqv"><img src="https://avatars.githubusercontent.com/u/23372024?v=4" title="tnqv" width="80" height="80"></a>
<a href="https://github.com/golemiso"><img src="https://avatars.githubusercontent.com/u/3282656?v=4" title="golemiso" width="80" height="80"></a>
<a href="https://github.com/sivchari"><img src="https://avatars.githubusercontent.com/u/55221074?v=4" title="sivchari" width="80" height="80"></a>
<a href="https://github.com/khanhtc3010"><img src="https://avatars.githubusercontent.com/u/9603918?v=4" title="khanhtc3010" width="80" height="80"></a>
<a href="https://github.com/p0tr3c"><img src="https://avatars.githubusercontent.com/u/12850042?v=4" title="p0tr3c" width="80" height="80"></a>
<a href="https://github.com/na-ga"><img src="https://avatars.githubusercontent.com/u/537006?v=4" title="na-ga" width="80" height="80"></a>
<a href="https://github.com/gotyoooo"><img src="https://avatars.githubusercontent.com/u/6133219?v=4" title="gotyoooo" width="80" height="80"></a>
<a href="https://github.com/ShotaKitazawa"><img src="https://avatars.githubusercontent.com/u/19530785?v=4" title="ShotaKitazawa" width="80" height="80"></a>
<a href="https://github.com/tennashi"><img src="https://avatars.githubusercontent.com/u/10219626?v=4" title="tennashi" width="80" height="80"></a>
<a href="https://github.com/Abirdcfly"><img src="https://avatars.githubusercontent.com/u/5100555?v=4" title="Abirdcfly" width="80" height="80"></a>
<a href="https://github.com/hongchaodeng"><img src="https://avatars.githubusercontent.com/u/920884?v=4" title="hongchaodeng" width="80" height="80"></a>
<a href="https://github.com/hori-ryota"><img src="https://avatars.githubusercontent.com/u/2936501?v=4" title="hori-ryota" width="80" height="80"></a>
<a href="https://github.com/eltociear"><img src="https://avatars.githubusercontent.com/u/22633385?v=4" title="eltociear" width="80" height="80"></a>
<a href="https://github.com/sano307"><img src="https://avatars.githubusercontent.com/u/12808316?v=4" title="sano307" width="80" height="80"></a>
<a href="https://github.com/misukuro"><img src="https://avatars.githubusercontent.com/u/1040546?v=4" title="misukuro" width="80" height="80"></a>
<a href="https://github.com/masaaania"><img src="https://avatars.githubusercontent.com/u/2755429?v=4" title="masaaania" width="80" height="80"></a>
<a href="https://github.com/KeisukeYamashita"><img src="https://avatars.githubusercontent.com/u/23056537?v=4" title="KeisukeYamashita" width="80" height="80"></a>
<a href="https://github.com/kentakozuka"><img src="https://avatars.githubusercontent.com/u/16733673?v=4" title="kentakozuka" width="80" height="80"></a>
<a href="https://github.com/Lennie"><img src="https://avatars.githubusercontent.com/u/330102?v=4" title="Lennie" width="80" height="80"></a>
<a href="https://github.com/kanata2"><img src="https://avatars.githubusercontent.com/u/7460883?v=4" title="kanata2" width="80" height="80"></a>
<a href="https://github.com/RikiyaFujii"><img src="https://avatars.githubusercontent.com/u/23261497?v=4" title="RikiyaFujii" width="80" height="80"></a>
<a href="https://github.com/SakataAtsuki"><img src="https://avatars.githubusercontent.com/u/58636635?v=4" title="SakataAtsuki" width="80" height="80"></a>
<a href="https://github.com/butterv"><img src="https://avatars.githubusercontent.com/u/15773082?v=4" title="butterv" width="80" height="80"></a>
<a href="https://github.com/mura-s"><img src="https://avatars.githubusercontent.com/u/4702673?v=4" title="mura-s" width="80" height="80"></a>
<a href="https://github.com/Linutux"><img src="https://avatars.githubusercontent.com/u/435352?v=4" title="Linutux" width="80" height="80"></a>
<a href="https://github.com/ww24"><img src="https://avatars.githubusercontent.com/u/695166?v=4" title="ww24" width="80" height="80"></a>
<a href="https://github.com/tnir"><img src="https://avatars.githubusercontent.com/u/10229505?v=4" title="tnir" width="80" height="80"></a>
<a href="https://github.com/yoiki"><img src="https://avatars.githubusercontent.com/u/39365493?v=4" title="yoiki" width="80" height="80"></a>
<a href="https://github.com/JohnTitor"><img src="https://avatars.githubusercontent.com/u/25030997?v=4" title="JohnTitor" width="80" height="80"></a>
<a href="https://github.com/mugioka"><img src="https://avatars.githubusercontent.com/u/62197019?v=4" title="mugioka" width="80" height="80"></a>
