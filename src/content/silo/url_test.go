package silo

import "testing"

func TestRedirectURL(t *testing.T) {
	link := "https://www.google.com/url?rct=j&sa=t&url=https://www.cryptoglobe.com/latest/2022/08/investment-strategist-lyn-alden-explains-why-she-is-still-bullish-on-bitcoin-long-term/&ct=ga&cd=CAIyGjlkMjI1NjUyODE3ODFjMDQ6Y29tOmVuOlVT&usg=AOvVaw16C2fJtw6m8QVEbto2HCKK"
	want := "https://www.cryptoglobe.com/latest/2022/08/investment-strategist-lyn-alden-explains-why-she-is-still-bullish-on-bitcoin-long-term/"
	have := RedirectURL(link)
	if have != want {
		t.Logf("want: %s", want)
		t.Logf("have: %s", have)
		t.Fail()
	}

	link = "https://example.com"
	if RedirectURL(link) != link {
		t.Fail()
	}

	link = "https://example.com/url?url=test.com"
	if RedirectURL(link) != link {
		t.Fail()
	}
}
