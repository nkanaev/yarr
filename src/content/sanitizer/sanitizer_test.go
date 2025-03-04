// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sanitizer

import "testing"

func TestValidInput(t *testing.T) {
	input := `<p>This is a <strong>text</strong> with an image: <img src="http://example.org/" alt="Test" loading="lazy">.</p>`
	want := `<p>This is a <strong>text</strong> with an image: <img src="http://example.org/" alt="Test" loading="lazy" referrerpolicy="no-referrer">.</p>`
	have := Sanitize("http://example.org/", input)

	if have != want {
		t.Errorf("Wrong output: \nwant: %#v\nhave: %#v", want, have)
	}
}

func TestImgWithTextDataURL(t *testing.T) {
	input := `<img src="data:text/plain;base64,SGVsbG8sIFdvcmxkIQ==" alt="Example">`
	expected := ``
	output := Sanitize("http://example.org/", input)

	if output != expected {
		t.Errorf(`Wrong output: %s`, output)
	}
}

func TestImgWithDataURL(t *testing.T) {
	input := `<img src="data:image/gif;base64,test" alt="Example">`
	want := `<img src="data:image/gif;base64,test" alt="Example" loading="lazy" referrerpolicy="no-referrer">`
	have := Sanitize("http://example.org/", input)

	if have != want {
		t.Errorf("Wrong output:\nwant: %s\nhave: %s", want, have)
	}
}

func TestImgWithSrcset(t *testing.T) {
	input := `<img srcset="example-320w.jpg, example-480w.jpg 1.5x,   example-640w.jpg 2x, example-640w.jpg 640w" src="example-640w.jpg" alt="Example">`
	want := `<img srcset="http://example.org/example-320w.jpg, http://example.org/example-480w.jpg 1.5x, http://example.org/example-640w.jpg 2x, http://example.org/example-640w.jpg 640w" src="http://example.org/example-640w.jpg" alt="Example" loading="lazy" referrerpolicy="no-referrer">`
	have := Sanitize("http://example.org/", input)

	if have != want {
		t.Errorf("Wrong output:\nwant: %s\nhave: %s", want, have)
	}
}

func TestImgWithSrcsetAndDataURL(t *testing.T) {
	input := `<img srcset="data:image/gif;base64,test" src="http://example.org/example-320w.jpg" alt="Example">`
	want := `<img srcset="data:image/gif;base64,test" src="http://example.org/example-320w.jpg" alt="Example" loading="lazy" referrerpolicy="no-referrer">`
	have := Sanitize("http://example.org/", input)

	if have != want {
		t.Errorf("Wrong output:\nwant: %s\nhave: %s", want, have)
	}
}

func TestSourceWithSrcsetAndMedia(t *testing.T) {
	input := `<picture><source media="(min-width: 800px)" srcset="elva-800w.jpg"></picture>`
	expected := `<picture><source media="(min-width: 800px)" srcset="http://example.org/elva-800w.jpg"></picture>`
	output := Sanitize("http://example.org/", input)

	if output != expected {
		t.Errorf(`Wrong output: %s`, output)
	}
}

func TestMediumImgWithSrcset(t *testing.T) {
	input := `<img alt="Image for post" class="t u v ef aj" src="https://miro.medium.com/max/5460/1*aJ9JibWDqO81qMfNtqgqrw.jpeg" srcset="https://miro.medium.com/max/552/1*aJ9JibWDqO81qMfNtqgqrw.jpeg 276w, https://miro.medium.com/max/1000/1*aJ9JibWDqO81qMfNtqgqrw.jpeg 500w" sizes="500px" width="2730" height="3407">`
	want := `<img alt="Image for post" src="https://miro.medium.com/max/5460/1*aJ9JibWDqO81qMfNtqgqrw.jpeg" srcset="https://miro.medium.com/max/552/1*aJ9JibWDqO81qMfNtqgqrw.jpeg 276w, https://miro.medium.com/max/1000/1*aJ9JibWDqO81qMfNtqgqrw.jpeg 500w" sizes="500px" loading="lazy" referrerpolicy="no-referrer">`
	have := Sanitize("http://example.org/", input)

	if have != want {
		t.Errorf("Wrong output:\nwant: %s\nhave: %s", want, have)
	}
}

func TestSelfClosingTags(t *testing.T) {
	input := `<p>This <br> is a <strong>text</strong><br/>.</p>`
	output := Sanitize("http://example.org/", input)

	if input != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, input, output)
	}
}

func TestTable(t *testing.T) {
	input := `<table><tr><th>A</th><th colspan="2">B</th></tr><tr><td>C</td><td>D</td><td>E</td></tr></table>`
	output := Sanitize("http://example.org/", input)

	if input != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, input, output)
	}
}

func TestRelativeURL(t *testing.T) {
	input := `This <a href="/test.html">link is relative</a> and this image: <img src="../folder/image.png"/>`
	want := `This <a href="http://example.org/test.html" rel="noopener noreferrer" target="_blank" referrerpolicy="no-referrer">link is relative</a> and this image: <img src="http://example.org/folder/image.png" loading="lazy" referrerpolicy="no-referrer"/>`
	have := Sanitize("http://example.org/", input)

	if want != have {
		t.Errorf("Wrong output:\nwant: %s\nhave: %s", want, have)
	}
}

func TestProtocolRelativeURL(t *testing.T) {
	input := `This <a href="//static.example.org/index.html">link is relative</a>.`
	expected := `This <a href="http://static.example.org/index.html" rel="noopener noreferrer" target="_blank" referrerpolicy="no-referrer">link is relative</a>.`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestInvalidTag(t *testing.T) {
	input := `<p>My invalid <wtf>tag</wtf>.</p>`
	expected := `<p>My invalid tag.</p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestVideoTag(t *testing.T) {
	input := `<p>My valid <video src="videofile.webm" autoplay poster="posterimage.jpg">fallback</video>.</p>`
	expected := `<p>My valid <video src="http://example.org/videofile.webm" poster="http://example.org/posterimage.jpg" controls>fallback</video>.</p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestAudioAndSourceTag(t *testing.T) {
	input := `<p>My music <audio controls="controls"><source src="foo.wav" type="audio/wav"></audio>.</p>`
	expected := `<p>My music <audio controls><source src="http://example.org/foo.wav" type="audio/wav"></audio>.</p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestUnknownTag(t *testing.T) {
	input := `<p>My invalid <unknown>tag</unknown>.</p>`
	expected := `<p>My invalid tag.</p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestInvalidNestedTag(t *testing.T) {
	input := `<p>My invalid <wtf>tag with some <em>valid</em> tag</wtf>.</p>`
	expected := `<p>My invalid tag with some <em>valid</em> tag.</p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestValidIFrame(t *testing.T) {
	input := `<iframe src="http://example.org/"></iframe>`
	want := `<iframe src="http://example.org/" sandbox="allow-scripts allow-same-origin allow-popups" loading="lazy"></iframe>`
	have := Sanitize("http://example.org/", input)

	if want != have {
		t.Errorf("Wrong output:\nwant: %s\nhave: %s", want, have)
	}
}

func TestInvalidIFrame(t *testing.T) {
	input := `<iframe src="http://example.org/"></iframe>`
	expected := ``
	output := Sanitize("http://example.com/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestIFrameWithChildElements(t *testing.T) {
	input := `<iframe src="https://www.youtube.com/"><p>test</p></iframe>`
	expected := `<div class="video-wrapper"><iframe src="https://www.youtube.com/" sandbox="allow-scripts allow-same-origin allow-popups" loading="lazy"></iframe></div>`
	output := Sanitize("http://example.com/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestInvalidURLScheme(t *testing.T) {
	input := `<p>This link is <a src="file:///etc/passwd">not valid</a></p>`
	expected := `<p>This link is not valid</p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestMailtoURIScheme(t *testing.T) {
	input := `<p>This link is <a href="mailto:jsmith@example.com?subject=A%20Test&amp;body=My%20idea%20is%3A%20%0A">valid</a></p>`
	expected := `<p>This link is <a href="mailto:jsmith@example.com?subject=A%20Test&amp;body=My%20idea%20is%3A%20%0A" rel="noopener noreferrer" target="_blank" referrerpolicy="no-referrer">valid</a></p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestTelURIScheme(t *testing.T) {
	input := `<p>This link is <a href="tel:+1-201-555-0123">valid</a></p>`
	expected := `<p>This link is <a href="tel:+1-201-555-0123" rel="noopener noreferrer" target="_blank" referrerpolicy="no-referrer">valid</a></p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestXMPPURIScheme(t *testing.T) {
	input := `<p>This link is <a href="xmpp:user@host?subscribe&amp;type=subscribed">valid</a></p>`
	expected := `<p>This link is <a href="xmpp:user@host?subscribe&amp;type=subscribed" rel="noopener noreferrer" target="_blank" referrerpolicy="no-referrer">valid</a></p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestBlacklistedLink(t *testing.T) {
	input := `<p>This image is not valid <img src="https://stats.wordpress.com/some-tracker"></p>`
	expected := `<p>This image is not valid </p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestXmlEntities(t *testing.T) {
	input := `<pre>echo "test" &gt; /etc/hosts</pre>`
	expected := `<pre>echo &#34;test&#34; &gt; /etc/hosts</pre>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestEspaceAttributes(t *testing.T) {
	input := `<td rowspan="<b>test</b>">test</td>`
	expected := `<td rowspan="&lt;b&gt;test&lt;/b&gt;">test</td>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestReplaceIframeURL(t *testing.T) {
	input := `<iframe src="https://player.vimeo.com/video/123456?title=0&amp;byline=0"></iframe>`
	expected := `<div class="video-wrapper"><iframe src="https://player.vimeo.com/video/123456?title=0&amp;byline=0" sandbox="allow-scripts allow-same-origin allow-popups" loading="lazy"></iframe></div>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestReplaceNoScript(t *testing.T) {
	input := `<p>Before paragraph.</p><noscript>Inside <code>noscript</code> tag with an image: <img src="http://example.org/" alt="Test" loading="lazy"></noscript><p>After paragraph.</p>`
	expected := `<p>Before paragraph.</p><p>After paragraph.</p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestReplaceScript(t *testing.T) {
	input := `<p>Before paragraph.</p><script type="text/javascript">alert("1");</script><p>After paragraph.</p>`
	expected := `<p>Before paragraph.</p><p>After paragraph.</p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestReplaceStyle(t *testing.T) {
	input := `<p>Before paragraph.</p><style>body { background-color: #ff0000; }</style><p>After paragraph.</p>`
	expected := `<p>Before paragraph.</p><p>After paragraph.</p>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf(`Wrong output: "%s" != "%s"`, expected, output)
	}
}

func TestWrapYoutubeIFrames(t *testing.T) {
	input := `<iframe src="https://www.youtube.com/embed/foobar"></iframe>`
	expected := `<div class="video-wrapper"><iframe src="https://www.youtube.com/embed/foobar" sandbox="allow-scripts allow-same-origin allow-popups" loading="lazy"></iframe></div>`
	output := Sanitize("http://example.org/", input)

	if expected != output {
		t.Errorf("Wrong output:\nwant: %v\nhave: %v", expected, output)
	}
}
