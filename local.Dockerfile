FROM scratch
COPY dist/cnfuzz /
EXPOSE 80
ENTRYPOINT ["/cnfuzz"]
