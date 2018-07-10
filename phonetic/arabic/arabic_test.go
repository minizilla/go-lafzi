package arabic

import (
	"testing"
)

func TestNormalizedUthmani(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		{[]byte(string([]rune{MaddahA, AlefA, SHLigatureSad, SHLigatureQaf,
			SHMeemInit, SHLamAlef, AlefWasla, SHJeem, SHThreeDots,
			SHSeen, RubElHizb, SHURectZero, SWaw, SHMeemIsolated,
			SLSeen, Sajdah, ECHStop, HamzaA, RHFCStop, SLMeem, Tatweel})),
			string([]rune{Alef, Hamza})},
		{[]byte("اقْتَرَبَ"), "إِقْتَرَبَ"},
		{[]byte("اقْرَ"), "إِقْرَ"},
	}

	for _, table := range tables {
		actual := string(normalizedUthmani(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestRemoveSpace(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Fatihah(1) verse: 1
		{[]byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ"), "بِسْمِاللَّهِالرَّحْمَنِالرَّحِيمِ"},
		// Yasin(36) verse: 2
		{[]byte("وَٱلْقُرْءَانِ ٱلْحَكِيمِ"), "وَالْقُرْءَانِالْحَكِيمِ"},
		// Al-Baqarah(2) verse: 249
		{[]byte("وَٱللَّهُ مَعَ ٱلصَّـٰبِرِينَ"), "وَاللَّهُمَعَالصَّبِرِينَ"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		actual := string(removeSpace(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestRemoveShadda(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Fatihah(1) verse: 1
		{[]byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ"), "بِسْمِاللَهِالرَحْمَنِالرَحِيمِ"},
		// Yasin(36) verse: 2
		{[]byte("وَٱلْقُرْءَانِ ٱلْحَكِيمِ"), "وَالْقُرْءَانِالْحَكِيمِ"},
		// Al-Baqarah(2) verse: 249
		{[]byte("وَٱللَّهُ مَعَ ٱلصَّـٰبِرِينَ"), "وَاللَهُمَعَالصَبِرِينَ"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		actual := string(removeShadda(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestJoinConsonant(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Kahfi(19) verse: 16 (idgham mutamatsilain: meem + sukun + meem -> meem)
		{[]byte("لَكُم مِّنْ أَمْرِكُم مِّرْفَقًا"), "لَكُمِنْأَمْرِكُمِرْفَقًا"},
		// Al-A'raf(7) verse: 160
		{[]byte("ٱضْرِب بِّعَصَاكَ ٱلْحَجَر"), "اضْرِبِعَصَاكَالْحَجَر"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		table.s = removeShadda(table.s)
		actual := string(joinConsonant(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestFixBoundary(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Fatihah(1) verse: 1 (harakat ending -> sukun)
		{[]byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ"), "بِسْمِالَهِالرَحْمَنِالرَحِيمْ"},
		// Al-Kahfi(19) verse: 16 (alef ending then fathatan -> fatha)
		{[]byte("لَكُم مِّنْ أَمْرِكُم مِّرْفَقًا"), "لَكُمِنْأَمْرِكُمِرْفَقَ"},
		// Al-A'raf(7) verse: 160 (consonant ending -> same, alef start -> alef hamza above + fatha)
		{[]byte("ٱضْرِب بِّعَصَاكَ ٱلْحَجَر"), "أَضْرِبِعَصَاكَالْحَجَر"},
		// Ar-Rum(30) verse: 21 (harakat ending then teh marbuta -> heh + sukun)
		{[]byte("مَّوَدَّةً وَرَحْمَةً"), "مَوَدَةًوَرَحْمَهْ"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		table.s = removeShadda(table.s)
		table.s = joinConsonant(table.s)
		actual := string(fixBoundary(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestTanwinSub(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Ar-Rum(30) verse: 21 (fathatan -> fatha + noon + sukun)
		{[]byte("مَّوَدَّةً وَرَحْمَةً"), "مَوَدَةَنْوَرَحْمَهْ"},
		// Al-Baqarah(2) verse: 25 (kasratan -> kasra + noon + sukun)
		{[]byte("جَنَّـٰتٍ تَجْرِ"), "جَنَتِنْتَجْرْ"},
		// Al-Baqarah(2) verse: 143 (dammatan -> damma + noon + sukun)
		{[]byte("لَرَءُوفٌ رَّحِيمٌ"), "لَرَءُوفُنْرَحِيمْ"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		table.s = removeShadda(table.s)
		table.s = joinConsonant(table.s)
		table.s = fixBoundary(table.s)
		actual := string(tanwinSub(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestRemoveMadda(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Baqarah(2) verse: 143 (fatha + alef + non-harakat -> fatha + non-harakat)
		{[]byte("عَلَى ٱلنَّاسِ"), "عَلَالنَسْ"},
		// Ali-Imran(3) verse: 44 (kasra + yeh + non-harakat -> kasra + non-harakat)
		{[]byte("غَلِيظَ"), "غَلِظْ"},
		// Al-Baqarah(2) verse: 143 (damma + waw + non-harakat -> damma + non-harakat)
		{[]byte("لَرَءُوفٌ رَّحِيمٌ"), "لَرَءُفُنْرَحِمْ"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		table.s = removeShadda(table.s)
		table.s = joinConsonant(table.s)
		table.s = fixBoundary(table.s)
		table.s = tanwinSub(table.s)
		actual := string(removeMadda(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestRemoveUnreadConsonant(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Baqarah(2) verse: 143 (double remove: alef & lam)
		{[]byte("عَلَى ٱلنَّاسِ"), "عَلَنَسْ"},
		// Ali-Imran(3) verse: 44 (single remove: yeh)
		{[]byte("غَلِيظَ"), "غَلِظْ"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		table.s = removeShadda(table.s)
		table.s = joinConsonant(table.s)
		table.s = fixBoundary(table.s)
		table.s = tanwinSub(table.s)
		table.s = removeMadda(table.s)
		actual := string(removeUnreadConsonant(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestIqlabSub(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Baqarah(2) verse: 253 (noon + sukun + beh -> meem + sukun + beh)
		{[]byte("مِنۢ بَعْدِهِم"), "مِمْبَعْدِهِم"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		table.s = removeShadda(table.s)
		table.s = joinConsonant(table.s)
		table.s = fixBoundary(table.s)
		table.s = tanwinSub(table.s)
		table.s = removeMadda(table.s)
		table.s = removeUnreadConsonant(table.s)
		actual := string(iqlabSub(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestIdghamSub(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Ali-Imran(3) verse: 145 (exception: dunya -> same)
		{[]byte("دُّنْيَا نُۭ"), "دُنْيَنْ"},
		// As-Saff(61) verse: 4 (exception: bunyan -> same)
		{[]byte("بُنْيَـٰنٌۭ مَّر"), "بُنْيَنُمَر"},
		// Ar-Ra'd(13) verse: 4 (exception: sinwan -> same)
		{[]byte("صِنْوَانٌۭ"), "صِنْوَنْ"},
		// Al-An'am(6) verse: 99 (exception: qinwan -> same)
		{[]byte("قِنْوَانٌۭ"), "قِنْوَنْ"},
		// Al-Qalam(68) verse: 1 (exception: nunwalqalam -> same)
		{[]byte("نٌ ۚ وَٱلْقَلَمِ وَمَا يَسْطُرُونَ"), "نُنْوَلْقَلَمِوَمَيَسْطُرُنْ"},
		// Al-Baqarah(2) verse: 2 (idgham bighunnah: noon + sukun + meem -> meem)
		// 						  (idgham bilaghunnah: noon + sukun + reh -> reh)
		{[]byte("هُدًۭى مِّن رَّبِّهِمْ ۖ"), "هُدَمِرَبِهِمْ"},
		// Al-Baqarah(2) verse: 7 (idgham bighunnah: noon + sukun + waw -> waw)
		{[]byte("غِشَـٰوَةٌۭ ۖ وَلَهُم"), "غِشَوَةُوَلَهُم"},
		// Al-Baqarah(2) verse: 8, (idgham bighunnah: noon + sukun + yeh -> yeh)
		{[]byte("مَن يَقُولُ"), "مَيَقُلْ"},
		// Al-Baqarah(2) verse: 12, (idgham bilaghunnah: noon + sukun + lam -> lam)
		{[]byte("وَلَـٰكِن لَّا يَشْعُرُونَ"), "وَلَكِلَيَشْعُرُنْ"},
		// Al-Baqarah(2) verse: 48 (idgham bighunnah: noon + sukun + noon -> noon)
		// note: noon + sukun + noon may be filtered out in joinConsonant
		{[]byte("عَنْ نَّفْسٍۢ"), "عَنَفْسْ"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		table.s = removeShadda(table.s)
		table.s = joinConsonant(table.s)
		table.s = fixBoundary(table.s)
		table.s = tanwinSub(table.s)
		table.s = removeMadda(table.s)
		table.s = removeUnreadConsonant(table.s)
		table.s = iqlabSub(table.s)
		actual := string(idghamSub(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestRemoveHarakat(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Ar-Rum(30) verse: 21 (fathatan -> fatha + noon + sukun)
		{[]byte("مَّوَدَّةً وَرَحْمَةً"), "مودةورحمه"},
		// Al-Baqarah(2) verse: 25 (kasratan -> kasra + noon + sukun)
		{[]byte("جَنَّـٰتٍ تَجْرِ"), "جنتنتجر"},
		// Al-Baqarah(2) verse: 143 (dammatan -> damma + noon + sukun)
		{[]byte("لَرَءُوفٌ رَّحِيمٌ"), "لرءفرحم"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		table.s = removeShadda(table.s)
		table.s = joinConsonant(table.s)
		table.s = fixBoundary(table.s)
		table.s = tanwinSub(table.s)
		table.s = removeMadda(table.s)
		table.s = removeUnreadConsonant(table.s)
		table.s = iqlabSub(table.s)
		table.s = idghamSub(table.s)
		actual := string(removeHarakat(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}

func TestEncode(t *testing.T) {
	tables := []struct {
		s        []byte
		expected string
	}{
		// Al-Fatihah(1) verse: 1-7 non-vowel
		{[]byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ"), "BSMLHRHMNRHM"},
		{[]byte("ٱلْحَمْدُ لِلَّهِ رَبِّ ٱلْعَـٰلَمِين"), "XLHMDLLHRBLXLMN"},
		{[]byte("ٱلرَّحْمَـٰنِ ٱلرَّحِيم"), "XRHMNRHM"},
		{[]byte("مَـٰلِكِ يَوْمِ ٱلدِّين"), "MLKYWMDN"},
		{[]byte("إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ"), "XYKNXBDWXYKNSTXN"},
		{[]byte("ٱهْدِنَا ٱلصِّرَٰطَ ٱلْمُسْتَقِيمَ"), "XHDNSRTLMSTKM"},
		{[]byte("صِرَٰطَ ٱلَّذِينَ أَنْعَمْتَ عَلَيْهِمْ غَيْرِ ٱلْمَغْضُوبِ عَلَيْهِمْ وَلَا ٱلضَّآلِّينَ"), "SRTLZNXNXMTXLYHMGYRLMGDBXLYHMWLDLN"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		table.s = removeShadda(table.s)
		table.s = joinConsonant(table.s)
		table.s = fixBoundary(table.s)
		table.s = tanwinSub(table.s)
		table.s = removeMadda(table.s)
		table.s = removeUnreadConsonant(table.s)
		table.s = iqlabSub(table.s)
		table.s = idghamSub(table.s)
		table.s = removeHarakat(table.s)
		actual := string(encode(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}

	tables = []struct {
		s        []byte
		expected string
	}{
		// Al-Fatihah(1) verse: 1-7 with vowel
		{[]byte("بِسْمِ ٱللَّهِ ٱلرَّحْمَـٰنِ ٱلرَّحِيمِ"), "BISMILAHIRAHMANIRAHIM"},
		{[]byte("ٱلْحَمْدُ لِلَّهِ رَبِّ ٱلْعَـٰلَمِين"), "XALHAMDULILAHIRABILXALAMIN"},
		{[]byte("ٱلرَّحْمَـٰنِ ٱلرَّحِيم"), "XARAHMANIRAHIM"},
		{[]byte("مَـٰلِكِ يَوْمِ ٱلدِّين"), "MALIKIYAWMIDIN"},
		{[]byte("إِيَّاكَ نَعْبُدُ وَإِيَّاكَ نَسْتَعِينُ"), "XIYAKANAXBUDUWAXIYAKANASTAXIN"},
		{[]byte("ٱهْدِنَا ٱلصِّرَٰطَ ٱلْمُسْتَقِيمَ"), "XAHDINASIRATALMUSTAKIM"},
		{[]byte("صِرَٰطَ ٱلَّذِينَ أَنْعَمْتَ عَلَيْهِمْ غَيْرِ ٱلْمَغْضُوبِ عَلَيْهِمْ وَلَا ٱلضَّآلِّينَ"), "SIRATALAZINAXANXAMTAXALAYHIMGAYRILMAGDUBIXALAYHIMWALADALIN"},
	}

	for _, table := range tables {
		table.s = normalizedUthmani(table.s)
		table.s = removeSpace(table.s)
		table.s = removeShadda(table.s)
		table.s = joinConsonant(table.s)
		table.s = fixBoundary(table.s)
		table.s = tanwinSub(table.s)
		table.s = removeMadda(table.s)
		table.s = removeUnreadConsonant(table.s)
		table.s = iqlabSub(table.s)
		table.s = idghamSub(table.s)
		actual := string(encode(table.s)[:])
		if actual != table.expected {
			t.Errorf("expected: %s, actual: %s", table.expected, actual)
		}
	}
}
