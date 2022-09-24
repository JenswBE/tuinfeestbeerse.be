let carouselClock = document.getElementById("carousel-clock");
let carouselClockDays = document.getElementById("carousel-clock-days");
let carouselClockHours = document.getElementById("carousel-clock-hours");
let carouselClockMinutes = document.getElementById("carousel-clock-minutes");
let carouselClockSeconds = document.getElementById("carousel-clock-seconds");
let carouselClockOngoing = document.getElementById("carousel-clock-ongoing");
let carouselClockAfter = document.getElementById("carousel-clock-after");
const start = new Date(
  Date.parse(carouselClock.getAttribute("data-start-date"))
);
const end = new Date(Date.parse(carouselClock.getAttribute("data-end-date")));
const SEC = 1000;
const MIN = 60 * SEC;
const HRS = 60 * MIN;
const DAY = 24 * HRS;

function updateClock() {
  const now = new Date();
  if (now < start) {
    // Start is in future => Show clock
    // Based on https://stackoverflow.com/a/59793084
    const diff = start - now;
    const days = Math.floor(diff / DAY);
    const hours = Math.floor((diff % DAY) / HRS);
    const mins = Math.floor((diff % HRS) / MIN);
    const secs = Math.floor((diff % MIN) / SEC);

    carouselClockDays.textContent = days;
    carouselClockHours.textContent = hours;
    carouselClockMinutes.textContent = mins;
    carouselClockSeconds.textContent = secs;
  } else if (now < end) {
    // Party ongoing
    carouselClock.id = "carousel-clock-hidden";
    carouselClock.classList.remove("d-md-flex");
    carouselClockOngoing.classList.remove("d-none");
  } else {
    // After party
    carouselClock.id = "carousel-clock-hidden";
    carouselClock.classList.remove("d-md-flex");
    carouselClockAfter.classList.remove("d-none");
  }
}

updateClock();
setInterval(updateClock, 1000);
