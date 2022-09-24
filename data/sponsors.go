package data

func getSponsors() Sponsors {
	return Sponsors{
		Show: false,
		MainSponsors: []*Sponsor{
			{
				Name:    "Campine",
				Link:    "https://www.campine.com/",
				Picture: "campine.svg",
			},
			{
				Name:    "Vastgoedpartners",
				Link:    "https://www.vastgoedpartners.be/nl",
				Picture: "vastgoed-small.jpg",
			},
		},
		Sponsors: []*Sponsor{
			{
				Name:    "Effenaf Graaf!",
				Link:    "https://www.effenafgraaf.be/",
				Picture: "effenaf-graaf.png",
			},
			{
				Name:    "Klaverblad",
				Link:    "https://www.facebook.com/drankenhandelklaverblad",
				Picture: "klaverblad.jpg",
			},
			{
				Name:    "GR Technics",
				Link:    "http://www.gr-technics.be",
				Picture: "gr.png",
			},
			{
				Name:    "Argenta - Kantoor Wouters",
				Link:    "https://www.argenta.be/nl/kantoren/advieskantoor-wouters-bvba-3552.html",
				Picture: "argenta.jpg",
			},
			{
				Name:    "Wesa Productions",
				Link:    "http://www.wesaproductions.be",
				Picture: "wesa.png",
			},
			{
				Name:    "Hermans",
				Link:    "https://www.hermans-heftrucks.be",
				Picture: "hermans.png",
			},
			{
				Name:    "Uit in Beerse",
				Link:    "http://www.uitinbeerse.be/",
				Picture: "uit-in-beerse.png",
			},
			{
				Name:    "CM",
				Link:    "https://cm.be/",
				Picture: "cm.jpg",
			},
			{
				Name:    "'t Broojke",
				Link:    "http://tbroojke.be",
				Picture: "broojke.jpg",
			},
		},
	}
}

type Sponsors struct {
	// Show the sponsors section
	Show bool
	// List of main sponsors. These are shown bigger than the other sponsors.
	// Sponsors are shown in same order as this list.
	MainSponsors []*Sponsor
	// List of sponsors. Sponsors are shown in same order.
	Sponsors []*Sponsor
}

type Sponsor struct {
	// Name of the sponsor
	Name string
	// Link to the website of the sponsor
	Link string
	// Name of the sponsor picture. E.g. klaverblad.jpg
	Picture string
}
