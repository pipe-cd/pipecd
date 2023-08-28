# https://github.com/pipe-cd/pipecd/pkgs/container/piped-base/103341355?tag=v0.44.0-alpha-52-ge881cac
FROM 713346901986.dkr.ecr.ap-northeast-1.amazonaws.com/azurite-piped/piped-base:20230825172803

ADD .artifacts/launcher /usr/local/bin/launcher
RUN mkdir -p /home/piped/.cache/piped-launcher/piped/bin
RUN chmod -R 755 /home/piped/.cache/piped-launcher/piped/bin 
ADD .artifacts/piped_v0.45.0_linux_amd64 /home/piped/.cache/piped-launcher/piped/bin/piped

ENTRYPOINT ["launcher"]
