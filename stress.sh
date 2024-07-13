vegeta attack -targets=targets.txt -rate=10000 -duration=60s | tee results.bin | vegeta report
vegeta report -type=text results.bin > report.txt
vegeta report -type=html results.bin > report.html
