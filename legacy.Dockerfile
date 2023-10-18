FROM gcr.io/distroless/static
COPY vulcain /
CMD ["/vulcain"]
EXPOSE 80 443

