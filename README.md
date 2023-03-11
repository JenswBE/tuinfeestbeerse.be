# Tuinfeest Beerse

```bash
# Preprocess data
# Re-run this step if files changed in "data" folder.
go run .

# Install Hugo
# See https://gohugo.io/installation/linux
go install github.com/gohugoio/hugo@latest

# Start Hugo server
hugo server --source website/

# Production build
hugo --source website/ --destination ../public
```

## Directory structure

- `data`: Contains all non-technical data of the website like artists, timetable, ...
- `website/content`: Contains all pages of the website EXCEPT the home page
- `website/layouts/_default/baseof.html`: Contains the base layout of the website
- `website/layouts/partials`: Contains reusable snippets of code. Snippets specific to a single page are in their own folder, e.g. `index`
- `website/layouts/index.html`: Layout of the home page
- `website/static`: Files which will be copied to the root of the site on build
