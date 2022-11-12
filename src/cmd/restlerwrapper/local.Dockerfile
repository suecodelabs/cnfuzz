FROM mcr.microsoft.com/restlerfuzzer/restler:v9.1.0 as final
COPY dist/restlerwrapper /
COPY src/cmd/restlerwrapper/auth.py /scripts/auth.py
ENTRYPOINT ["/restlerwrapper"]
