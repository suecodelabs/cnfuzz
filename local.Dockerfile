FROM scratch
COPY dist/cnfuzz /
EXPOSE 8080
CMD ["/cnfuzz"]
