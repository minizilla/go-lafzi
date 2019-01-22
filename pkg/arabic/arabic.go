package arabic

// Consonant characters.
var (
	Hamza       = '\u0621' // http://www.fileformat.info/info/unicode/char/0621/index.htm
	AlefMaddaA  = '\u0622' // http://www.fileformat.info/info/unicode/char/0622/index.htm
	AlefHamzaA  = '\u0623' // http://www.fileformat.info/info/unicode/char/0623/index.htm
	WawHamzaA   = '\u0624' // http://www.fileformat.info/info/unicode/char/0624/index.htm
	AlefHamzaB  = '\u0625' // http://www.fileformat.info/info/unicode/char/0625/index.htm
	YehHamzaA   = '\u0626' // http://www.fileformat.info/info/unicode/char/0626/index.htm
	Alef        = '\u0627' // http://www.fileformat.info/info/unicode/char/0627/index.htm
	Beh         = '\u0628' // http://www.fileformat.info/info/unicode/char/0628/index.htm
	TehMarbuta  = '\u0629' // http://www.fileformat.info/info/unicode/char/0629/index.htm
	Teh         = '\u062a' // http://www.fileformat.info/info/unicode/char/062a/index.htm
	Theh        = '\u062b' // http://www.fileformat.info/info/unicode/char/062b/index.htm
	Jeem        = '\u062c' // http://www.fileformat.info/info/unicode/char/062c/index.htm
	Hah         = '\u062d' // http://www.fileformat.info/info/unicode/char/062d/index.htm
	Khah        = '\u062e' // http://www.fileformat.info/info/unicode/char/062e/index.htm
	Dal         = '\u062f' // http://www.fileformat.info/info/unicode/char/062f/index.htm
	Thal        = '\u0630' // http://www.fileformat.info/info/unicode/char/0630/index.htm
	Reh         = '\u0631' // http://www.fileformat.info/info/unicode/char/0631/index.htm
	Zain        = '\u0632' // http://www.fileformat.info/info/unicode/char/0632/index.htm
	Seen        = '\u0633' // http://www.fileformat.info/info/unicode/char/0633/index.htm
	Sheen       = '\u0634' // http://www.fileformat.info/info/unicode/char/0634/index.htm
	Sad         = '\u0635' // http://www.fileformat.info/info/unicode/char/0635/index.htm
	Dad         = '\u0636' // http://www.fileformat.info/info/unicode/char/0636/index.htm
	Tah         = '\u0637' // http://www.fileformat.info/info/unicode/char/0637/index.htm
	Zah         = '\u0638' // http://www.fileformat.info/info/unicode/char/0638/index.htm
	Ain         = '\u0639' // http://www.fileformat.info/info/unicode/char/0639/index.htm
	Ghain       = '\u063a' // http://www.fileformat.info/info/unicode/char/063a/index.htm
	Feh         = '\u0641' // http://www.fileformat.info/info/unicode/char/0641/index.htm
	Qaf         = '\u0642' // http://www.fileformat.info/info/unicode/char/0642/index.htm
	Kaf         = '\u0643' // http://www.fileformat.info/info/unicode/char/0643/index.htm
	Lam         = '\u0644' // http://www.fileformat.info/info/unicode/char/0644/index.htm
	Meem        = '\u0645' // http://www.fileformat.info/info/unicode/char/0645/index.htm
	Noon        = '\u0646' // http://www.fileformat.info/info/unicode/char/0646/index.htm
	Heh         = '\u0647' // http://www.fileformat.info/info/unicode/char/0647/index.htm
	Waw         = '\u0648' // http://www.fileformat.info/info/unicode/char/0648/index.htm
	AlefMaksura = '\u0649' // http://www.fileformat.info/info/unicode/char/0649/index.htm
	Yeh         = '\u064a' // http://www.fileformat.info/info/unicode/char/064a/index.htm
)

// Vowels characters.
var (
	Fathatan = '\u064b' // http://www.fileformat.info/info/unicode/char/064b/index.htm
	Dammatan = '\u064c' // http://www.fileformat.info/info/unicode/char/064c/index.htm
	Kasratan = '\u064d' // http://www.fileformat.info/info/unicode/char/064d/index.htm
	Fatha    = '\u064e' // http://www.fileformat.info/info/unicode/char/064e/index.htm
	Damma    = '\u064f' // http://www.fileformat.info/info/unicode/char/064f/index.htm
	Kasra    = '\u0650' // http://www.fileformat.info/info/unicode/char/0650/index.htm
	Shadda   = '\u0651' // http://www.fileformat.info/info/unicode/char/0651/index.htm
	Sukun    = '\u0652' // http://www.fileformat.info/info/unicode/char/0652/index.htm
)

// Uthmani characters. Prefix: S = small, H = high, L = low, U = upright,
// E = empty, C = centre, R = Rounded, F = filled
var (
	Tatweel        = '\u0640' // http://www.fileformat.info/info/unicode/char/0640/index.htm
	MaddahA        = '\u0653' // http://www.fileformat.info/info/unicode/char/0653/index.htm
	HamzaA         = '\u0654' // http://www.fileformat.info/info/unicode/char/0654/index.htm
	AlefA          = '\u0670' // http://www.fileformat.info/info/unicode/char/0670/index.htm
	AlefWasla      = '\u0671' // http://www.fileformat.info/info/unicode/char/0671/index.htm
	SHLigatureSad  = '\u06d6' // http://www.fileformat.info/info/unicode/char/06d6/index.htm
	SHLigatureQaf  = '\u06d7' // http://www.fileformat.info/info/unicode/char/06d7/index.htm
	SHMeemInit     = '\u06d8' // http://www.fileformat.info/info/unicode/char/06d8/index.htm
	SHLamAlef      = '\u06d9' // http://www.fileformat.info/info/unicode/char/06d9/index.htm
	SHJeem         = '\u06da' // http://www.fileformat.info/info/unicode/char/06da/index.htm
	SHThreeDots    = '\u06db' // http://www.fileformat.info/info/unicode/char/06db/index.htm
	SHSeen         = '\u06dc' // http://www.fileformat.info/info/unicode/char/06dc/index.htm
	RubElHizb      = '\u06de' // http://www.fileformat.info/info/unicode/char/06de/index.htm
	SHRZero        = '\u06df' // http://www.fileformat.info/info/unicode/char/06df/index.htm
	SHURectZero    = '\u06e0' // http://www.fileformat.info/info/unicode/char/06e0/index.htm
	SHMeemIsolated = '\u06e2' // http://www.fileformat.info/info/unicode/char/06e2/index.htm
	SLSeen         = '\u06e3' // http://www.fileformat.info/info/unicode/char/06e3/index.htm
	SWaw           = '\u06e5' // http://www.fileformat.info/info/unicode/char/06e5/index.htm
	SYeh           = '\u06e6' // http://www.fileformat.info/info/unicode/char/06e6/index.htm
	SHYeh          = '\u06e7' // http://www.fileformat.info/info/unicode/char/06e7/index.htm
	SHNoon         = '\u06e8' // http://www.fileformat.info/info/unicode/char/06e8/index.htm
	Sajdah         = '\u06e9' // http://www.fileformat.info/info/unicode/char/06e9/index.htm
	ECLStop        = '\u06ea' // http://www.fileformat.info/info/unicode/char/06ea/index.htm
	ECHStop        = '\u06eb' // http://www.fileformat.info/info/unicode/char/06eb/index.htm
	RHFCStop       = '\u06ec' // http://www.fileformat.info/info/unicode/char/06ec/index.htm
	SLMeem         = '\u06ed' // http://www.fileformat.info/info/unicode/char/06ed/index.htm
)
