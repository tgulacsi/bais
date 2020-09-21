/*
Package bais provides functions to Encode and Decode a more compact
alternative Base64 encoding

Encodes and decodes BAIS (Byte Array In String) encoding,
which preserves runs of ASCII characters unchanged. This encoding is
useful for debugging (since ASCII runs are visible) and for conversion
of bytes to JSON.</summary>

Arrays encoded with <see cref="ByteArrayInString.Convert(ArraySlice{byte}, bool)"/>
tend to be slightly more compact than standard Uuencoding or Base64,
and when you use this encoding in JSON with UTF-8, the output is
typically also more compact than yEnc since double-byte characters
above 127 are avoided.

A BAIS string alternates between runs of "direct" bytes (usually bytes
in the ASCII range that are represented as themselves) and runs of a
special base-64 encoding. The base-64 encoding is a sequence of 6-bit
digits with 64 added to them, except for 63 which is mapped to itself.
This is easier and faster to encode and decode than standard Base64
and has an interesting property described below.

A BAIS string begins in ASCII mode and switches to base 64 when the '\b'
character is encountered. Base-64 mode ends, returning to ASCII, when a
'!' character is encountered.

For example:
<pre>
b := byte[]{ 67, 97, 116, 128, 10, 69, 255, 65, 66, 67, 68 };
bytes.Equal(bais.Encode(b), "Cat\b`@iE?tEB!CD"), b)
"Cat\b`@iE?tEB!CD" = bais.Decode(bais.Encode(b))
</pre>
A byte sequence such as 128, 10, 69, 255 can be encoded in base 64 as
illustrated:
<pre>
           ---128---    ---10----    ---69----  ---255---
Bytes:     1000 0000    0000 1010    0100 0101  1111 1111
Base 64:   100000   000000   101001    000101   111111   110000
Encoded: 01100000 01000000 01101001  01000101 01111111 01110000
         ---96--- ---64--- --105---  ---69--- --127--- --112---
            `        @        i         E        ~        p
</pre>

An interesting property of this base-64 encoding is that when it encodes
bytes between 63 and 126, those bytes appear unchanged at certain
offsets (specifically the third, sixth, ninth, etc.) In this example,
since the third byte is 'E' (69), it also appears as 'E' in the
output.

When viewing BAIS strings, another thing to keep in mind is that
runs of zeroes ('\0') will tend to appear as runs of `@` characters
in the base 64 encoding, although a single zero is not always enough
to make a `@` appear. Runs of 255 will tend to appear as runs of `?`.

There are many ways to encode a given byte array as BAIS.
*/
package bais
