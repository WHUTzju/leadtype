package codepage

var CP1253 = Ranges{
	{0x0000, 0x007F, 128, 0},
	{0x00A0, 0x00A0, 1, 0},
	{0x00A3, 0x00A9, 7, 0},
	{0x00AB, 0x00AE, 4, 0},
	{0x00B0, 0x00B3, 4, 0},
	{0x00B5, 0x00B7, 3, 0},
	{0x00BB, 0x00BB, 1, 0},
	{0x00BD, 0x00BD, 1, 0},
	{0x0192, 0x0192, 1, -271},
	{0x0384, 0x0384, 1, -720},
	{0x0385, 0x0386, 2, -740},
	{0x0388, 0x038A, 3, -720},
	{0x038C, 0x038C, 1, -720},
	{0x038E, 0x03A1, 20, -720},
	{0x03A3, 0x03CE, 44, -720},
	{0x2013, 0x2014, 2, -8061},
	{0x2015, 0x2015, 1, -8038},
	{0x2018, 0x2019, 2, -8071},
	{0x201A, 0x201A, 1, -8088},
	{0x201C, 0x201D, 2, -8073},
	{0x201E, 0x201E, 1, -8090},
	{0x2020, 0x2021, 2, -8090},
	{0x2022, 0x2022, 1, -8077},
	{0x2026, 0x2026, 1, -8097},
	{0x2030, 0x2030, 1, -8103},
	{0x2039, 0x2039, 1, -8110},
	{0x203A, 0x203A, 1, -8095},
	{0x20AC, 0x20AC, 1, -8236},
	{0x2122, 0x2122, 1, -8329},
}
