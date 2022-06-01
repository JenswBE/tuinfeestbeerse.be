let carouselClock = document.getElementById("carousel-clock");
let carouselClockDays = document.getElementById("carousel-clock-days");
let carouselClockHours = document.getElementById("carousel-clock-hours");
let carouselClockMinutes = document.getElementById("carousel-clock-minutes");
let carouselClockSeconds = document.getElementById("carousel-clock-seconds");
const start = Date.parse(carouselClock.getAttribute("data-date"));
const SEC = 1000;
const MIN = 60 * SEC;
const HRS = 60 * MIN;
const DAY = 24 * HRS;

function updateClock() {
  const diff = start - new Date();

  // Based on https://stackoverflow.com/a/59793084
  const days = Math.floor(diff / DAY);
  const hours = Math.floor((diff % DAY) / HRS);
  const mins = Math.floor((diff % HRS) / MIN);
  const secs = Math.floor((diff % MIN) / SEC);

  carouselClockDays.textContent = days;
  carouselClockHours.textContent = hours;
  carouselClockMinutes.textContent = mins;
  carouselClockSeconds.textContent = secs;
}

updateClock();
setInterval(updateClock, 1000);
