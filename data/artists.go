package data

import "strings"

func getArtists() Artists {
	return Artists{
		Show: true,
		Artists: []*Artist{
			{
				Name:    "Sub-lime",
				Type:    "Rock band",
				Picture: "sub-lime.jpg",
				Description: strings.Join([]string{
					"Sub-Lime is een 5-koppige rockband met roots in de Kempen.",
					"Met een setlist geïnspireerd op iconische bands zoals Kiss, QOTSA, Foo Fighters, ...",
					"krijgen ze jong en oud samen op de dansvloer.",
				}, " "),
				Links: map[string]string{
					"globe":    "https://www.sub-lime.be/",
					"facebook": "https://www.facebook.com/watch/sublimeband/",
				},
			},
			{
				Name:    "The Skadillacs",
				Type:    "Ska band",
				Picture: "skadillacs.jpg",
				Links: map[string]string{
					"globe":    "http://www.theskadillacs.be/",
					"facebook": "https://www.facebook.com/theskadillacs",
				},
			},
			{
				Name:    "De Romeo's",
				Type:    "Ambiance band",
				Picture: "romeos.jpg",
				Links: map[string]string{
					"globe":     "https://deromeos.be/",
					"facebook":  "https://www.facebook.com/deromeos",
					"twitter":   "https://twitter.com/romeosofficial",
					"instagram": "https://www.instagram.com/deromeos/",
				},
			},
			{
				Name:    "De Koperdieven",
				Type:    "Koper cover band",
				Picture: "koperdieven.jpg",
				Description: strings.Join([]string{
					"Een zootje losgeslagen ongeregeld die niets liever doen dan elk feestje op zijn kop te zetten.",
					"Dit doen ze door dansbare, meezingbare en brulbare nummers op een bijna strafbare manier op jullie los te laten.",
					"De koperdieven&hellip; ze stelen de show!",
				}, " "),
				Links: map[string]string{
					"facebook": "https://www.facebook.com/koperdieven",
				},
			},
			{
				Name:    "Winnaar DJ contest @ JH Den Beir",
				Picture: "den-beir.jpg",
				Description: strings.Join([]string{
					"De enige echte winnaar van de",
					`<a href="https://www.facebook.com/events/729280071543041" style="color: lightblue; font-weight: bold">DJ contest van JH Den Beir</a>.`,
				}, " "),
				Links: map[string]string{
					"facebook": "https://www.facebook.com/events/729280071543041",
				},
			},
			{
				Name:    "DJ N1CO",
				Type:    "All-round DJ",
				Picture: "n1co.jpg",
				Links: map[string]string{
					"facebook": "https://www.facebook.com/N1COfficial",
				},
			},
			{
				Name:    "DJ Lenny",
				Type:    "All-round DJ",
				Picture: "lenny.jpg",
				Links: map[string]string{
					"facebook": "https://www.facebook.com/DJLenny007/",
				},
			},
			{
				Name:    "DJ Double-U",
				Type:    "All-round DJ",
				Picture: "double-u.jpg",
				Links: map[string]string{
					"globe":      "https://www.djdouble-u.be/",
					"facebook":   "https://www.facebook.com/Double.U.DJ/",
					"instagram":  "https://www.instagram.com/double_u_dj/",
					"soundcloud": "https://soundcloud.com/dj-doubleu",
				},
			},
			{
				Name:    "DJ Polleke & Celis",
				Type:    "All-round DJ's",
				Picture: "polleke-celis.jpg",
				Description: strings.Join([]string{
					"Wat is er beter dan DJ Polleke? 2 Pollekes natuurlijk!<br>",
					"Dit jaar slaan Polleke en zijn zoon Celis, de handen in elkaar om u nog vettere schijven,",
					"zottere dansplaten en coolere hits te brengen. Een combi dus om U tegen te zeggen en dé moment",
					"om nog eens volledig los te gaan!",
				}, " "),
				Links: map[string]string{
					"facebook":   "https://www.facebook.com/DJ-POLLEKE-225905228432/",
					"youtube":    "https://www.youtube.com/channel/UCwE_Gdx71HxVaOxERTX9wgg",
					"facebook-2": "https://www.facebook.com/celisofficial",
					"soundcloud": "https://soundcloud.com/celisbe",
				},
			},
			{
				Name:    "Vzw Salsmanians",
				Type:    "Salsa",
				Picture: "salsa-manians.jpg",
				Description: strings.Join([]string{
					"Opgericht op 1 juli 2009. Dansclub van vrijwilligers, voor salsa &amp; andere caribische dansen.",
					"Wij verzorgen danslessen, workshops, initiaties, party's, demo's, salsadj-sets, ...",
				}, " "),
				Links: map[string]string{
					// "globe":    "https://www.salsmanians.be/",
					"facebook": "https://www.facebook.com/salsmanians/",
				},
			},
		},
	}
}

type Artists struct {
	// Show the artists section
	Show bool
	// List of artists. Artists are shown in same order.
	Artists []*Artist
}

type Artist struct {
	// Name of the artist
	Name string
	// Type of artist. E.g. Rock band, All-round DJ, ...
	Type string
	// Name of the artist picture. E.g. lenny.jpg
	Picture string
	// Description of the artist
	Description string
	// Map of links for the artist.
	// Key is icon name from https://icons.getbootstrap.com/.
	// Value is URL link should point to.
	Links map[string]string
}
