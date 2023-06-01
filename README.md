# Tuinfeest Beerse

```bash
# Generate website
go run .

# Local development
# => Regenerate website on change
go install github.com/cespare/reflex@latest
reflex -vR 'output/.*' -- go run . --debug
```
