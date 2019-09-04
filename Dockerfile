FROM scratch
COPY vulcain /
COPY public ./public/
CMD ["./vulcain"]
EXPOSE 80 443
