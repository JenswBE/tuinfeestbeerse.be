# Optimize images

```bash
cd website/static/assets/images/carousels/main
rm *.jpg
cp ../../carousels-original/main/*.jpg .
for img in *.jpg
do
    convert "${img:?}" -resize 1500x1000 "${img:?}"
done
jpegoptim --strip-all --all-progressive --max 75 *.jpg
```
