(function($) {
    "use strict";

  // Responsive-menu trigger
    $(".menu").on('click', function() {
        $(".responsive-menu-area").toggleClass("active");
    });

    $('ul.metismenu').metisMenu({
    });

    // Slider-active
    $('.slider-active').owlCarousel({
        margin: 0,
        loop: true,
        nav: false,
        autoplayTimeout: 4000,
        smartSpeed: 1200,
        autoplay: true,
        items: 1
    });

    // Slider-active
    $(".slider-active").on('translate.owl.carousel', function() {
        $('.slider-items h2').removeClass('fadeInLeft animated').hide();
        $('.slider-items p').removeClass('fadeInUp animated').hide();
        $('.slider-items ul').removeClass('fadeInUp animated').hide();
    });

    $(".slider-active").on('translated.owl.carousel', function() {
        $('.owl-item.active .slider-items h2').addClass('fadeInLeft animated').show();
        $('.owl-item.active .slider-items p').addClass('fadeInUp animated').show();
        $('.owl-item.active .slider-items ul').addClass('fadeInUp animated').show();
    });

    //slider-area background setting
    function sliderBgSetting() {
        if ($(".slider-active .slider-items,.slider-active .slider-items").length) {
            $(".slider-active .slider-items,.slider-active .slider-items").each(function() {
                var $this = $(this);
                var img = $this.find(".slider").attr("src");

                $this.css({
                    backgroundImage: "url(" + img + ")",
                    backgroundSize: "cover",
                    backgroundPosition: "center center"
                })
            });
        }
    }
    sliderBgSetting()

    // Sticky menu
    $(window).on('scroll', function() {
        var scroll = $(window).scrollTop(),
            mainHeader = $('#sticky-header'),
            mainHeaderHeight = mainHeader.innerHeight();

        // console.log(mainHeader.innerHeight());
        if (scroll > 0) {
            $("#sticky-header").addClass("sticky-menu");
        } else {
            $("#sticky-header").removeClass("sticky-menu");
        }
    });

    /*--------------------------
     scrollUp
    ---------------------------- */
    $.scrollUp({
        scrollText: '<i class="fa fa-arrow-up"></i>',
        easingType: 'linear',
        scrollSpeed: 900,
        animation: 'fade'
    });

    function setTwoColEqHeight($col1, $col2) {
        var firstCol = $col1,
            secondCol = $col2,
            firstColHeight = $col1.innerHeight(),
            secondColHeight = $col2.innerHeight();

        if (firstColHeight > secondColHeight) {
            secondCol.css({
                "height": firstColHeight + 1 + "px"
            })
        } else {
            firstCol.css({
                "height": secondColHeight + 1 + "px"
            })
        }
    }

    // Smooth-scrolling
    function smoothScrolling($links, $topGap) {
        var links = $links;
        var topGap = $topGap;

        links.on("click", function() {
            if (location.pathname.replace(/^\//, '') === this.pathname.replace(/^\//, '') && location.hostname === this.hostname) {
                var target = $(this.hash);
                target = target.length ? target : $("[name=" + this.hash.slice(1) + "]");
                if (target.length) {
                    $("html, body").animate({
                        scrollTop: target.offset().top - topGap
                    }, 1000, "easeInOutExpo");
                    return false;
                }
            }
            return false;
        });
    }

    $(window).on("load", function() {
        smoothScrolling($(".smooth-links a[href^='#']"), -2);
    });


    /*---------------------
     Countdown
    --------------------- */
    $('[data-countdown]').each(function() {
        var $this = $(this),
            finalDate = $(this).data('countdown');
        $this.countdown(finalDate, function(event) {
            $this.html(event.strftime('<span class="cdown days"><span class="time-count">%-D</span> <p>Dagen</p></span> <span class="cdown hour"><span class="time-count">%-H</span> <p>Uren</p></span> <span class="cdown minutes"><span class="time-count">%M</span> <p>Minuten</p></span> <span class="cdown second"> <span><span class="time-count">%S</span> <p>Seconden</p></span>'));
        });
    });

})(jQuery);