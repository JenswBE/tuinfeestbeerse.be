<%
# Imports
from sys import exit
from os import listdir
from random import shuffle
from datetime import date, datetime
import logging

# Remap data to local vars (support file renaming)
general_conf = data['general']
sponsors_conf = data['sponsors']
football_conf = data['football']

# Dates
dt_fb_start = datetime.strptime(football_conf['start'], '%d.%m.%Y %H:%M')
dt_fb_end = datetime.strptime(football_conf['end'], '%d.%m.%Y %H:%M')
dt_fb_start2 = datetime.strptime(football_conf['start2'], '%d.%m.%Y %H:%M')
dt_fb_end2 = datetime.strptime(football_conf['end2'], '%d.%m.%Y %H:%M')

# Error helper
def exit_with_error(message, code=1):
  logging.error(message)
  exit(1)

# Get backgrounds
backgrounds_url = 'assets/images/football/slider'
backgrounds = listdir('static/' + backgrounds_url)
backgrounds.sort()

# Get other vars
name_facebook = general_conf['links'].get('facebook')['link'].strip('/').split('/')[-1] if general_conf['links'].get('facebook') else ''

# Define helpers
def link(url, content, title=False, blank=True):
  target = 'target="{}"'.format('_blank' if blank else '_self')
  title = 'title="{}"'.format(title) if title else ''
  return "<a href=\"{}\" {} {}>{}</a>".format(url, target, title, content)

def link_icon(url, icon, title=False, appendix=False, blank=True):
  icon = "<i class=\"fa fa-{}\"></i>".format(icon)
  content = "{} {}".format(icon, appendix) if appendix else icon
  return link(url, content, title, blank)
%>

<%def name="menu_links()">
  <li><a href="/">Tuinfeest ${dt_fb_start.year}</a></li>
  %if sponsors_conf['settings']['show_sponsors']:
      <li><a href="#partners">Sponsors</a></li>
  %endif
  <li><a href="#map">Waar?</a></li>
</%def>

<!doctype html>
<html class="no-js" lang="nl">
  <head>
      <meta charset="utf-8">
      <meta http-equiv="x-ua-compatible" content="ie=edge">
      <title>Tuinfeest Beerse</title>
      <meta name="description" content="">
      <meta name="viewport" content="width=device-width, initial-scale=1">

      <!-- Favicon -->
      <link rel="apple-touch-icon" sizes="57x57" href="/assets/favicon/apple-icon-57x57.png">
      <link rel="apple-touch-icon" sizes="60x60" href="/assets/favicon/apple-icon-60x60.png">
      <link rel="apple-touch-icon" sizes="72x72" href="/assets/favicon/apple-icon-72x72.png">
      <link rel="apple-touch-icon" sizes="76x76" href="/assets/favicon/apple-icon-76x76.png">
      <link rel="apple-touch-icon" sizes="114x114" href="/assets/favicon/apple-icon-114x114.png">
      <link rel="apple-touch-icon" sizes="120x120" href="/assets/favicon/apple-icon-120x120.png">
      <link rel="apple-touch-icon" sizes="144x144" href="/assets/favicon/apple-icon-144x144.png">
      <link rel="apple-touch-icon" sizes="152x152" href="/assets/favicon/apple-icon-152x152.png">
      <link rel="apple-touch-icon" sizes="180x180" href="/assets/favicon/apple-icon-180x180.png">
      <link rel="icon" type="image/png" sizes="192x192"  href="/assets/favicon/android-icon-192x192.png">
      <link rel="icon" type="image/png" sizes="32x32" href="/assets/favicon/favicon-32x32.png">
      <link rel="icon" type="image/png" sizes="96x96" href="/assets/favicon/favicon-96x96.png">
      <link rel="icon" type="image/png" sizes="16x16" href="/assets/favicon/favicon-16x16.png">
      <link rel="manifest" href="/assets/favicon/manifest.json">
      <meta name="msapplication-TileColor" content="#ffffff">
      <meta name="msapplication-TileImage" content="/assets/favicon/ms-icon-144x144.png">
      <meta name="theme-color" content="#ffffff">

      <!-- Stylesheets -->
      <link rel="stylesheet" href="/assets/css/bootstrap.min.css">
      <link rel="stylesheet" href="/assets/css/animate.css">
      <link rel="stylesheet" href="/assets/css/owl.carousel.min.css">
      <link rel="stylesheet" href="/assets/css/swiper.min.css">
      <link rel="stylesheet" href="/assets/css/font-awesome.min.css">
      <link rel="stylesheet" href="/assets/css/flaticon.css">
      <link rel="stylesheet" href="/assets/css/magnific-popup.css">
      <link rel="stylesheet" href="/assets/css/metisMenu.min.css">
      <link rel="stylesheet" href="/assets/css/styles.css">
      <link rel="stylesheet" href="/assets/css/responsive.css">

      <!-- Modernizr -->
      <script src="/assets/js/vendor/modernizr-2.8.3.min.js"></script>

      <!-- Leaflet Maps API -->
      <link rel="stylesheet" href="https://unpkg.com/leaflet@1.4.0/dist/leaflet.css" integrity="sha512-puBpdR0798OZvTTbP4A8Ix/l+A4dHDD0DGqYW6RQ+9jxkRFclaxxQb/SJAWZfWAkuyeQUytO7+7N4QKrDh+drA==" crossorigin=""/>
      <script src="https://unpkg.com/leaflet@1.4.0/dist/leaflet.js" integrity="sha512-QVftwZFqvtRNi0ZyCtsznlKSWOStnDORoefr1enyq5mVL4tmKB3S/EnC3rRJcxCPavG10IcrVGSmPh6Qw5lwrg==" crossorigin=""></script>
  </head>
  <body>
      <!--[if lt IE 8]>
              <p class="browserupgrade">You are using an <strong>outdated</strong> browser. Please <a href="http://browsehappy.com/">upgrade your browser</a> to improve your experience.</p>
          <![endif]-->
      <!-- header-area start -->
      <header class="header-area">
          <div class="container-fluid">
              <div class="row">
                  <div class="col-lg-4 col-md-8 col-8">
                      <div class="logo">
                          <a href="/">
                              <img src="/assets/images/logo.png" alt="">
                          </a>
                      </div>
                  </div>
                  <div class="col-lg-8 d-none d-lg-block">
                      <div class="mainmenu">
                          <ul class="d-flex justify-content-end smooth-links">
                              ${menu_links()}
                          </ul>
                      </div>
                  </div>
                  <div class="d-block d-lg-none col-md-4 col-4 pull-right col-sm-4">
                      <ul class="menu">
                          <li class="first"></li>
                          <li class="second"></li>
                          <li class="third"></li>
                      </ul>
                  </div>
              </div>
          </div>
          <!-- responsive-menu area start -->
          <div class="responsive-menu-area d-block d-lg-none">
              <div class="container">
                  <div class="row">
                      <div class="col-12">
                          <ul class="metismenu smooth-links">
                              ${menu_links()}
                          </ul>
                      </div>
                  </div>
              </div>
          </div>
          <!-- responsive-menu area start -->
      </header>
      <!-- header-area end -->
      <!-- slider-area start -->
      <div class="slider-area">
          <div class="slider-active owl-carousel">
            %for bg in backgrounds:
              <div class="slider-items">
                  <img src="/${backgrounds_url}/${bg}" alt="" class="slider">
              </div>
            %endfor
          </div>
          <div id="slider-static" class="slider-content">
                <div class="col-lg-10 col-12 offset-lg-1 text-center">
                    <h2>EK Finland - Belgi&euml;</h2>
                    <p>${football_conf["day_name"].capitalize()} ${dt_fb_start.day} Juni ${dt_fb_start.year}, ${dt_fb_start.strftime('%H:%M')} - ${dt_fb_end.strftime('%H:%M')}</p>
                    <div id="clock" data-countdown="${dt_fb_start.strftime('%Y/%m/%d %H:%M')}"></div>
                </div>
            </div>
          <ul class="social-share">
              <li>${link_icon("mailto:" + general_conf['mail-info'], "envelope", appendix="<span>{}</span>".format(general_conf['mail-info']))}</li>
              <li>${link_icon("mailto:" + general_conf['mail-artist'], "music", appendix="<span>{}</span>".format(general_conf['mail-artist']))}</li>
              %if name_facebook:
                <li>${link_icon(general_conf['links']['facebook'], "facebook", appendix="<span>{}</span>".format(name_facebook))}</li>
              %endif
          </ul>
      </div>
      <!-- slider-area end -->
      <!-- about-area start -->
      <div class="about-area" id="about">
          <div class="container">
              <div class="row">
                  <div class="col-lg-12 col-12">
                            %if football_conf['description_readmore']:
                      <div class="about-wrap readmore">
                      %else:
                      <div class="about-wrap">
                      %endif
                          ${football_conf['description']}
                      </div>
                  </div>
              </div>
          </div>
      </div>
      <!-- about-area end -->

    %if sponsors_conf['settings']['show_sponsors']:
        <div class="partner-area" id="partners">
            <div class="container">
                <div class="row">
                    <div class="col-12">
                        <div class="section-title text-center">
                            <h2>Onze Sponsors</h2>
                            <p>Bedankt voor jullie steun aan Tuinfeest ${dt_fb_start.year}!</p>
                        </div>
                    </div>
                </div>

                <div class="row">
                    <div class="col-lg-6 offset-lg-3 col-sm-8 offset-sm-2 col-12">
                      <div class="row">
                        %for sponsor in sponsors_conf['main-sponsors']:
                              <div class="col-12">
                                  <div class="partner-wrap">
                                      <a href="${sponsor['link']}" target="_blank">
                                          <img src="/assets/images/sponsor/${sponsor['picture']}" alt="${sponsor['title']}">
                                      </a>
                                  </div>
                              </div>
                          %endfor
                        </div>
                    </div>
                </div>

                <div class="row">
                    <div class="col-lg-8 offset-lg-2 col-sm-10 offset-sm-1 col-12">
                      <div class="row">
                        %for sponsor in sponsors_conf['sponsors']:
                              <div class="col-sm-4 col-6">
                                  <div class="partner-wrap">
                                      <a href="${sponsor['link']}" target="_blank">
                                          <img src="/assets/images/sponsor/${sponsor['picture']}" alt="${sponsor['title']}">
                                      </a>
                                  </div>
                              </div>
                          %endfor
                        </div>
                    </div>
                </div>
            </div>
        </div>
    %endif

      <!-- .content-area start -->
      <div class="content-area-thin">
          <div class="container">
          </div>
      </div>
      <!-- .content-area end -->

      <div id="map" hight="400"></div>
      <!-- footer-area start -->
      <footer class="footer-area">
          <div class="container">
              <div class="row">
                  <div class="col-md-8 col-12">
                      <div class="copyright">
                          <p>&copy; <a href="https://jensw.be" target="_blank">Jensw.be</a> ${date.today().year}</p>
                      </div>
                  </div>
                  <div class="col-md-4 col-12">
                      <ul class="social-icon justify-content-end d-flex">
                          %for name, item in general_conf['links'].items():
                              %if name:
                                  <li>${link_icon(item['link'], name, item['title'])}</li>
                              %endif
                          %endfor
                      </ul>
                  </div>
              </div>
          </div>
      </footer>
      <!-- footer-area end -->
      <!-- jquery latest version -->
      <script src="/assets/js/vendor/jquery-2.2.4.min.js"></script>
      <!-- popper.min.js -->
      <script src="/assets/js/vendor/popper.min.js"></script>
      <!-- bootstrap js -->
      <script src="/assets/js/bootstrap.min.js"></script>
      <!-- owl.carousel.2.0.0-beta.2.4 css -->
      <script src="/assets/js/owl.carousel.min.js"></script>
      <!-- swiper.min.js -->
      <script src="/assets/js/swiper.min.js"></script>
      <!-- plugins js -->
      <script src="/assets/js/jquery.canvasjs.min.js"></script>
      <!-- metisMenu.min.js -->
      <script src="/assets/js/metisMenu.min.js"></script>
      <!-- plugins js -->
      <script src="/assets/js/plugins.js"></script>
      <!-- Leaflet Map -->
      <script>
    jQuery(document).ready(function(){
        // Basic setup
        var coord_nief_park = [51.31681, 4.85721];
        var mymap = L.map('map', {
          dragging: false,
          scrollWheelZoom: false,
          inertia: false,
          noMoveStart: true,
          tap: false
        }).setView(coord_nief_park, 16);

        // Add map tile layer
        L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token={accessToken}', {
          attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, <a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
          maxZoom: 18,
          id: 'mapbox.streets',
          accessToken: 'pk.eyJ1IjoidHVpbmZlZXN0YmVlcnMiLCJhIjoiY2pxeTU1YzQxMDAxZzQ1cGV5NGlieGpnbyJ9.frYTXarGgo6JlWsXrtLs9A'
        }).addTo(mymap);

        var tfIcon = L.icon({
          iconUrl: '/assets/images/icons/leaflet/tf-marker.png',
          iconRetinaUrl: '/assets/images/icons/leaflet/tf-marker-2x.png',
          iconSize:     [75, 75], // size of the icon
          iconAnchor:   [37, 37],
          popupAnchor:  [0, -38]
      });

        // Add marker for Nief Park
        L.Icon.Default.prototype.options.imagePath = '/assets/images/icons/leaflet/';
        var marker = L.marker(coord_nief_park, {icon: tfIcon}).addTo(mymap);
        marker.bindPopup("NIEF PARK<br>Ingang via:<br>Pastoriestraat 8<br>2340 Beerse").openPopup();
    });
    </script>
      <!-- main js -->
      <script src="/assets/js/scripts.js"></script>
  </body>
</html>
