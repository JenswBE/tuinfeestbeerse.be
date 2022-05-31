// Basic setup
const coord_nief_park = {
  lat: 51.31681,
  lng: 4.85765,
};
var mymap = L.map("tf-map", {
  dragging: false,
  scrollWheelZoom: false,
  inertia: false,
  noMoveStart: true,
  tap: false,
}).setView(coord_nief_park, 16);

// Add map tile layer
L.tileLayer(
  "https://api.mapbox.com/styles/v1/{id}/tiles/{z}/{x}/{y}?access_token={accessToken}",
  {
    attribution:
      'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, <a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
    maxZoom: 18,
    tileSize: 512,
    zoomOffset: -1,
    id: "mapbox/streets-v11",
    accessToken:
      "pk.eyJ1IjoidHVpbmZlZXN0YmVlcnMiLCJhIjoiY2pxeTU1YzQxMDAxZzQ1cGV5NGlieGpnbyJ9.frYTXarGgo6JlWsXrtLs9A",
  }
).addTo(mymap);

var tfIcon = L.icon({
  iconUrl: "/assets/images/icons/leaflet/tf-marker-groen.png",
  iconRetinaUrl: "/assets/images/icons/leaflet/tf-marker-2x-groen.png",
  iconSize: [50, 50], // size of the icon
  iconAnchor: [25, 25],
  popupAnchor: [0, -25],
});

// Add marker for Nief Park
L.Icon.Default.prototype.options.imagePath = "/assets/images/icons/leaflet/";
var marker = L.marker(coord_nief_park, { icon: tfIcon }).addTo(mymap);
marker
  .bindPopup("NIEF PARK<br>Ingang via:<br>Pastoriestraat 8<br>2340 Beerse")
  .openPopup();
