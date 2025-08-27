package templates

func Revision20120111() Template {
	return &template{
		gridDefinition: map[int]RetrieveEndOctet{
			0:     constantEnd(72),
			1:     constantEnd(84),
			2:     constantEnd(84),
			3:     constantEnd(96),
			4:     constantEnd(48),
			5:     constantEnd(60),
			10:    constantEnd(72),
			12:    constantEnd(84),
			13:    constantEnd(88),
			20:    constantEnd(65),
			23:    constantEnd(81),
			30:    constantEnd(81),
			31:    constantEnd(81),
			33:    constantEnd(97),
			40:    constantEnd(72),
			41:    constantEnd(84),
			42:    constantEnd(84),
			43:    constantEnd(96),
			50:    constantEnd(28),
			51:    constantEnd(40),
			52:    constantEnd(40),
			53:    constantEnd(52),
			61:    constantEnd(112),
			62:    constantEnd(106),
			63:    constantEnd(121),
			90:    constantEnd(80),
			100:   constantEnd(38),
			101:   constantEnd(35),
			110:   constantEnd(57),
			140:   constantEnd(64),
			150:   constantEnd(42),
			204:   constantEnd(72),
			1100:  constantEnd(82),
			32768: constantEnd(72),
			32769: constantEnd(80),
		},
		productDefinition: map[int]RetrieveEndOctet{
			0: constantEnd(34),
			1: constantEnd(37),
			2: constantEnd(36),
			5: constantEnd(47),
			6: constantEnd(35),
			7: constantEnd(34),
			8: func(bytes []byte) int {
				timeRanges := int(bytes[41])
				end := 46 + (12 * timeRanges)
				return end
			},
			9: func(bytes []byte) int {
				timeRanges := int(bytes[54])
				end := 59 + (12 * timeRanges)
				return end
			},
		},
	}
}
