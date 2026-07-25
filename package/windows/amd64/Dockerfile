# escape=`
FROM mcr.microsoft.com/windows/nanoserver:ltsc2022

WORKDIR C:\pasturestack
COPY internal-dns.exe C:\pasturestack\internal-dns.exe
COPY start.cmd C:\pasturestack\start.cmd

USER ContainerAdministrator
ENTRYPOINT ["C:\\pasturestack\\start.cmd"]
