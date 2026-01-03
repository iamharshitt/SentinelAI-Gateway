\- name: Build SentinelAI Agent

&nbsp; working-directory: agent

&nbsp; run: |

&nbsp;   go mod tidy

&nbsp;   go build ./...



