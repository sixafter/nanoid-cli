// Copyright (c) 2024 Six After, Inc
//
// This source code is licensed under the Apache 2.0 License found in the
// LICENSE file in the root directory of this source tree.

package nanoid

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"math/bits"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"github.com/sixafter/nanoid/x/crypto/prng"
)

var (
	// DefaultGenerator is a global, shared instance of a Nano ID generator. It is safe for concurrent use.
	DefaultGenerator Generator

	// DefaultRandReader is the default random number generator used for generating IDs.
	DefaultRandReader = prng.Reader

	// ErrDuplicateCharacters is returned when the provided alphabet contains duplicate characters.
	ErrDuplicateCharacters = errors.New("duplicate characters in alphabet")

	// ErrExceededMaxAttempts is returned when the maximum number of attempts to perform
	// an operation, such as generating a unique ID, has been exceeded.
	ErrExceededMaxAttempts = errors.New("exceeded maximum attempts")

	// ErrInvalidLength is returned when a specified length value for an operation is invalid.
	ErrInvalidLength = errors.New("invalid length")

	// ErrInvalidAlphabet is returned when the provided alphabet for generating IDs is invalid.
	ErrInvalidAlphabet = errors.New("invalid alphabet")

	// ErrNonUTF8Alphabet is returned when the provided alphabet contains non-UTF-8 characters.
	ErrNonUTF8Alphabet = errors.New("alphabet contains invalid UTF-8 characters")

	// ErrAlphabetTooShort is returned when the provided alphabet has fewer than 2 characters.
	ErrAlphabetTooShort = errors.New("alphabet length is less than 2")

	// ErrAlphabetTooLong is returned when the provided alphabet exceeds 256 characters.
	ErrAlphabetTooLong = errors.New("alphabet length exceeds 256")

	// ErrNilRandReader is returned when the random number generator (rand.Reader) is nil,
	// preventing the generation of random values.
	ErrNilRandReader = errors.New("nil random reader")
)

const (
	// DefaultAlphabet defines the standard set of characters used for Nano ID generation.
	// It includes uppercase and lowercase English letters, digits, and the characters
	// '_' and '-'. This selection aligns with the Nano ID specification, ensuring
	// a URL-friendly and easily readable identifier.
	//
	// Example: "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	DefaultAlphabet = "_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// DefaultLength specifies the default number of characters in a generated Nano ID.
	// A length of 21 characters provides a high level of uniqueness while maintaining
	// brevity, making it suitable for most applications requiring unique identifiers.
	DefaultLength = 21

	// maxAttemptsMultiplier determines the maximum number of attempts the generator
	// will make to produce a valid Nano ID before failing. It is calculated as a
	// multiplier based on the desired ID length to balance between performance
	// and the probability of successful ID generation, especially when using
	// non-power-of-two alphabets.
	maxAttemptsMultiplier = 10

	// MinAlphabetLength sets the minimum permissible number of unique characters
	// in the alphabet used for Nano ID generation. An alphabet with fewer than
	// 2 characters would not provide sufficient variability for generating unique IDs,
	// making this a lower bound to ensure meaningful ID generation.
	//
	// Example: An alphabet like "AB" is acceptable, but "A" is not.
	MinAlphabetLength = 2

	// MaxAlphabetLength defines the maximum allowable number of unique characters
	// in the alphabet for Nano ID generation. This upper limit ensures that the
	// generator operates within reasonable memory and performance constraints,
	// preventing excessively large alphabets that could degrade performance or
	// complicate index calculations.
	MaxAlphabetLength = 256
)

// ID represents a Nano ID as a string.
type ID string

// EmptyID represents an empty Nano ID.
var EmptyID = ID("")

func init() {
	var err error
	DefaultGenerator, err = NewGenerator()
	if err != nil {
		panic(fmt.Sprintf("failed to initialize DefaultGenerator: %v", err))
	}
}

// ConfigOptions holds the configurable options for the Generator.
// It is used with the Function Options pattern.
type ConfigOptions struct {
	// RandReader is the source of randomness used for generating IDs.
	// By default, it uses x/crypto/prng/Reader, which provides cryptographically secure random bytes.
	RandReader io.Reader

	// Alphabet is the set of characters used to generate the Nano ID.
	// It must be a valid UTF-8 string containing between 2 and 256 unique characters.
	// Using a diverse and appropriately sized alphabet ensures the uniqueness and randomness of the generated IDs.
	Alphabet string

	// LengthHint specifies a typical or default length for generated IDs.
	LengthHint uint16
}

// Config holds the runtime configuration for the Nano ID generator.
//
// It is immutable after initialization and provides all the necessary
// parameters for generating unique IDs efficiently and securely.
type Config interface {
	// AlphabetLen returns the number of unique characters in the provided alphabet.
	//
	// This length determines the range of indices for selecting characters during ID generation.
	// Using uint16 allows for alphabets up to 65,535 characters.
	AlphabetLen() uint16

	// BaseMultiplier returns the foundational multiplier used in buffer size calculations.
	//
	// It is based on the logarithm of the intended ID length (LengthHint) plus 2.
	// This helps scale the buffer size appropriately with different ID lengths.
	BaseMultiplier() int

	// BitsNeeded returns the minimum number of bits required to represent all possible indices of the alphabet.
	//
	// This value is crucial for generating random numbers that map uniformly to the alphabet indices without bias.
	BitsNeeded() uint

	// BufferMultiplier returns the combined multiplier used in the buffer size calculation.
	//
	// It adds a fraction of the scaling factor to the base multiplier to fine-tune the buffer size,
	// considering both the ID length and the alphabet size.
	BufferMultiplier() int

	// BufferSize returns the total size of the buffer (in bytes) used for generating random data.
	//
	// The buffer size is calculated to balance efficiency and performance,
	// minimizing calls to the random number generator by reading larger chunks of random data at once.
	BufferSize() int

	// ByteAlphabet returns the slice of bytes representing the alphabet,
	// used when the alphabet consists solely of ASCII characters.
	//
	// For non-ASCII alphabets, this returns nil, and RuneAlphabet is used instead.
	ByteAlphabet() []byte

	// BytesNeeded returns the number of bytes required to store the BitsNeeded for each character in the ID.
	//
	// It rounds up BitsNeeded to the nearest byte, ensuring sufficient space for random data generation.
	BytesNeeded() uint

	// IsASCII returns true if the alphabet consists solely of ASCII characters.
	//
	// This allows for optimization in processing, using bytes instead of runes for ID generation.
	IsASCII() bool

	// IsPowerOfTwo returns true if the length of the alphabet is a power of two.
	//
	// When true, random index selection can be optimized using bitwise operations,
	// such as bitwise AND with the mask, improving performance.
	IsPowerOfTwo() bool

	// LengthHint returns the intended length of the IDs to be generated.
	//
	// This hint is used in calculations to adjust buffer sizes and scaling factors accordingly.
	LengthHint() uint16

	// MaxBytesPerRune represents the maximum number of bytes required to encode
	// any rune in the alphabet using UTF-8 encoding.
	//
	// This value is computed during
	// configuration based on the provided alphabet and is used to preallocate the
	// buffer size in the newUnicode function. By accurately estimating the buffer size,
	// we ensure efficient string building without unnecessary memory allocations
	// or buffer resizing.
	//
	// For example, if the alphabet includes only ASCII and Latin-1 characters, each rune
	// requires at most 2 bytes. However, if the alphabet includes emojis or other
	// multibyte characters, this value could be up to 4 bytes.
	MaxBytesPerRune() int

	// Mask returns the bitmask used to extract the necessary bits from randomly generated bytes.
	//
	// The mask is essential for efficiently mapping random values to valid alphabet indices,
	// ensuring uniform distribution and preventing bias.
	Mask() uint

	// RandReader returns the source of randomness used for generating IDs.
	//
	// It is typically a cryptographically secure random number generator (e.g., crypto/rand.Reader).
	RandReader() io.Reader

	// RuneAlphabet returns the slice of runes representing the alphabet.
	//
	// This is used for ID generation when the alphabet includes non-ASCII (multibyte) characters,
	// allowing support for a wider range of characters.
	RuneAlphabet() []rune

	// ScalingFactor returns the scaling factor used to adjust the buffer size.
	//
	// It balances the influence of the alphabet size and the intended ID length,
	// ensuring efficient random data generation without excessive memory usage.
	ScalingFactor() int
}

// Configuration defines the interface for retrieving generator configuration.
type Configuration interface {
	// Config returns the runtime configuration of the generator.
	Config() Config
}

// Generator defines the interface for generating Nano IDs.
// Implementations of this interface provide methods to create new IDs
// and to read random data, supporting both ID generation and direct random byte access.
type Generator interface {
	// New generates and returns a new Nano ID as a string with the specified length.
	// The 'length' parameter determines the number of characters in the generated ID.
	// Returns an error if the ID generation fails due to issues like insufficient randomness.
	//
	// Usage:
	//   id, err := generator.New(21)
	//   if err != nil {
	//       // handle error
	//   }
	//   fmt.Println("Generated ID:", id)
	New(length int) (ID, error)

	// Read fills the provided byte slice 'p' with random data, reading up to len(p) bytes.
	// Returns the number of bytes read and any error encountered during the read operation.
	//
	// Implements the io.Reader interface, allowing the Generator to be used wherever an io.Reader is accepted.
	// This can be useful for directly obtaining random bytes or integrating with other components that consume random data.
	//
	// Usage:
	//   buffer := make([]byte, 21)
	//   n, err := generator.Read(buffer)
	//   if err != nil {
	//       // handle error
	//   }
	//   fmt.Printf("Read %d random bytes\n", n)
	Read(p []byte) (n int, err error)
}

// runtimeConfig holds the runtime configuration for the Nano ID generator.
// It is immutable after initialization.
type runtimeConfig struct {
	// RandReader is the source of randomness used for generating IDs.
	randReader       io.Reader
	byteAlphabet     []byte
	runeAlphabet     []rune
	mask             uint
	bitsNeeded       uint
	bytesNeeded      uint
	bufferSize       int
	bufferMultiplier int
	scalingFactor    int
	baseMultiplier   int
	maxBytesPerRune  int
	alphabetLen      uint16
	lengthHint       uint16
	isASCII          bool
	isPowerOfTwo     bool
}

type generator struct {
	config      *runtimeConfig
	entropyPool *sync.Pool
	idPool      *sync.Pool
}

// New generates a new Nano ID using the default length specified by `DefaultLength`.
// It returns the generated ID as a string and any error encountered during the generation.
//
// Usage:
//
//	id, err := nanoid.New()
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Generated ID:", id)
func New() (ID, error) {
	return NewWithLength(DefaultLength)
}

// NewWithLength generates a new Nano ID of the specified length.
// It returns the generated ID as a string and any error encountered during the generation.
//
// Parameters:
//   - length int: The number of characters for the generated ID.
//
// Usage:
//
//	id, err := nanoid.NewWithLength(21)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Generated ID:", id)
func NewWithLength(length int) (ID, error) {
	return DefaultGenerator.New(length)
}

// Must generates a new Nano ID using the default length specified by `DefaultLength`.
// It returns the generated ID as a string.
// If an error occurs during ID generation, it panics.
// This function simplifies safe initialization of global variables holding pre-generated Nano IDs.
//
// Usage:
//
//	id := nanoid.Must()
//	fmt.Println("Generated ID:", id)
func Must() ID {
	return MustWithLength(DefaultLength)
}

// MustWithLength generates a new Nano ID of the specified length.
// It returns the generated ID as a string.
// If an error occurs during ID generation, it panics.
// The 'length' parameter specifies the number of characters in the generated ID.
// This function simplifies safe initialization of global variables holding pre-generated Nano IDs.
//
// Parameters:
//   - length int: The number of characters for the generated ID.
//
// Usage:
//
//	id := nanoid.MustWithLength(30)
//	fmt.Println("Generated ID:", id)
func MustWithLength(length int) ID {
	id, err := NewWithLength(length)
	if err != nil {
		panic(err)
	}

	return id
}

// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// Reader is the interface that wraps the basic Read method.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
//
// If len(p) == 0, Read should always return n == 0. It may return a
// non-nil error if some error condition is known, such as EOF.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
func Read(p []byte) (n int, err error) {
	return DefaultGenerator.Read(p)
}

// Option defines a function type for configuring the Generator.
// It allows for flexible and extensible configuration by applying
// various settings to the ConfigOptions during Generator initialization.
type Option func(*ConfigOptions)

// WithAlphabet sets a custom alphabet for the Generator.
// The provided alphabet string defines the set of characters that will be
// used to generate Nano IDs. This allows users to customize the character set
// according to their specific requirements, such as using only alphanumeric
// characters, including symbols, or supporting non-ASCII characters.
//
// Parameters:
//   - alphabet string: A string representing the desired set of characters for ID generation.
//
// Returns:
//   - Option: A configuration option that applies the custom alphabet to ConfigOptions.
//
// Usage:
//
//	generator, err := nanoid.NewGenerator(nanoid.WithAlphabet("abcdef123456"))
func WithAlphabet(alphabet string) Option {
	return func(c *ConfigOptions) {
		c.Alphabet = alphabet
	}
}

// WithRandReader sets a custom random reader for the Generator.
// By default, the Generator uses a cryptographically secure random number
// generator (e.g., crypto/rand.Reader). However, in some cases, users might
// want to provide their own source of randomness, such as for testing purposes
// or to integrate with a different entropy source.
//
// Parameters:
//   - reader io.Reader: An implementation of io.Reader that supplies random data.
//
// Returns:
//   - Option: A configuration option that applies the custom random reader to ConfigOptions.
//
// Usage Example:
//
//	 customReader := myCustomRandomReader()
//	 generator, err := nanoid.NewGenerator(
//		nanoid.WithRandReader(customReader))
func WithRandReader(reader io.Reader) Option {
	return func(c *ConfigOptions) {
		c.RandReader = reader
	}
}

// WithLengthHint sets the hint of the intended length of the IDs to be generated.
// Providing a length hint allows the Generator to optimize internal configurations,
// such as buffer sizes and scaling factors, based on the expected ID length. This
// can enhance performance and efficiency, especially when generating a large number
// of IDs with similar lengths.
//
// Parameters:
//   - hint uint16: A non-zero unsigned integer representing the anticipated length of the Nano IDs.
//
// Returns:
//   - Option: A configuration option that applies the length hint to ConfigOptions.
//
// Usage Example:
//
//	generator, err := nanoid.NewGenerator(nanoid.WithLengthHint(21))
func WithLengthHint(hint uint16) Option {
	return func(c *ConfigOptions) {
		c.LengthHint = hint
	}
}

// NewGenerator creates a new Generator with buffer pooling enabled.
// It accepts variadic Option parameters to configure the Generator's behavior.
// The function initializes the configuration with default values, applies any provided options,
// validates the configuration, constructs the runtime configuration, initializes buffer pools,
// and returns a configured Generator or an error if the configuration is invalid.
//
// Parameters:
//   - options ...Option: A variadic list of Option functions to customize the Generator's configuration.
//
// Returns:
//   - Generator: An instance of the Generator interface configured with the specified options.
//   - error: An error object if the Generator could not be created due to invalid configuration.
//
// Error Conditions:
//   - ErrInvalidLength: Returned if the provided LengthHint is less than 1.
//   - ErrNilRandReader: Returned if the provided RandReader is nil.
//   - ErrInvalidAlphabet: Returned if the alphabet is invalid or contains invalid UTF-8 characters.
//   - ErrNonUTF8Alphabet: Returned if the alphabet contains non-UTF-8 characters.
//   - ErrDuplicateCharacters: Returned if the alphabet contains duplicate characters.
func NewGenerator(options ...Option) (Generator, error) {
	// Initialize ConfigOptions with default values.
	// These defaults include the default alphabet, the default random reader,
	// and the default length hint for ID generation.
	configOpts := &ConfigOptions{
		Alphabet:   DefaultAlphabet,
		RandReader: DefaultRandReader,
		LengthHint: DefaultLength,
	}

	// Apply provided options to customize the configuration.
	// Each Option function modifies the ConfigOptions accordingly.
	for _, opt := range options {
		opt(configOpts)
	}

	// Ensure LengthHint is within valid bounds.
	// LengthHint must be at least 1 to generate meaningful IDs.
	if configOpts.LengthHint < 1 {
		return nil, ErrInvalidLength
	}

	// Ensure RandReader is not nil.
	// A valid randomness source is essential for generating secure IDs.
	if configOpts.RandReader == nil {
		return nil, ErrNilRandReader
	}

	// Validate and construct RuntimeConfig based on the current ConfigOptions.
	// buildRuntimeConfig performs validation on the alphabet and computes necessary
	// parameters for efficient ID generation.
	config, err := buildRuntimeConfig(configOpts)
	if err != nil {
		return nil, err
	}

	// Initialize a pool of byte slices for random data generation.
	// The pool helps in reusing memory buffers, reducing garbage collection overhead.
	entropyPool := &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, config.bufferSize*config.bufferMultiplier)
			return &buf
		},
	}

	var idPool *sync.Pool
	if config.isASCII {
		idPool = &sync.Pool{
			New: func() interface{} {
				buf := make([]byte, 0, config.bufferSize*config.bufferMultiplier)
				return &buf
			},
		}
	} else {
		idPool = &sync.Pool{
			New: func() interface{} {
				buf := make([]rune, 0, config.bufferSize*config.bufferMultiplier)
				return &buf
			},
		}
	}

	// Return the configured Generator instance.
	// The generator holds references to the runtime configuration and buffer pools,
	// facilitating efficient and thread-safe ID generation.
	return &generator{
		config:      config,
		entropyPool: entropyPool,
		idPool:      idPool,
	}, nil
}

func buildRuntimeConfig(opts *ConfigOptions) (*runtimeConfig, error) {
	if len(opts.Alphabet) == 0 {
		return nil, ErrInvalidAlphabet
	}

	// Check if the alphabet is valid UTF-8
	if !utf8.ValidString(opts.Alphabet) {
		return nil, ErrNonUTF8Alphabet
	}

	alphabetRunes := []rune(opts.Alphabet)
	isASCII := true
	byteAlphabet := make([]byte, len(alphabetRunes))
	maxBytesPerRune := 1 // Initialize to 1 for ASCII

	for i, r := range alphabetRunes {
		if r > unicode.MaxASCII {
			isASCII = false
			// Compute the number of bytes needed to encode this rune
			runeBytes := utf8.RuneLen(r)
			if runeBytes < 0 {
				return nil, ErrInvalidAlphabet
			}
			if runeBytes > maxBytesPerRune {
				maxBytesPerRune = runeBytes
			}
		} else {
			byteAlphabet[i] = byte(r)
		}
	}

	if !isASCII {
		// Convert to rune alphabet if non-ASCII characters are present
		byteAlphabet = nil // Clear byteAlphabet as it's not used
	}

	// Check for duplicate characters
	seenRunes := make(map[rune]bool)
	for _, r := range alphabetRunes {
		if seenRunes[r] {
			return nil, ErrDuplicateCharacters
		}
		seenRunes[r] = true
	}

	// The length of the alphabet, representing the number of unique characters available for ID generation.
	alphabetLen := uint16(len(alphabetRunes))

	// Ensure the alphabet length adheres to predefined constraints.
	if alphabetLen > MaxAlphabetLength {
		return nil, ErrAlphabetTooLong
	}

	if alphabetLen < MinAlphabetLength {
		return nil, ErrAlphabetTooShort
	}

	// Calculate the minimum number of bits needed to represent all indices of the alphabet.
	// This is essential for generating random numbers that map uniformly to the alphabet indices.
	// The calculation uses bits.Len to find the position of the highest set bit in alphabetLen - 1.
	bitsNeeded := uint(bits.Len(uint(alphabetLen - 1)))
	if bitsNeeded == 0 {
		return nil, ErrInvalidAlphabet
	}

	// Create a bitmask that isolates the bits needed to represent the alphabet indices.
	// The mask is used to efficiently extract valid bits from randomly generated bytes.
	mask := uint((1 << bitsNeeded) - 1)

	// TODO: Scale bitsNeeded based on length hint (???)
	//adjustedBitsNeeded := bitsNeeded + uint(math.Log2(float64(opts.LengthHint)))

	// Determine the number of bytes required to store 'bitsNeeded' bits, rounding up to the nearest byte.
	bytesNeeded := (bitsNeeded + 7) / 8

	// Check if the alphabet length is a power of two, allowing optimization of modulus operations using bitwise AND.
	// This optimization improves performance during random index generation.
	isPowerOfTwo := (alphabetLen & (alphabetLen - 1)) == 0

	// Calculate a base multiplier for buffer size based on the length hint.
	// The length hint indicates the desired length of the generated IDs.
	// Using logarithm ensures the buffer scales appropriately with the ID length.
	baseMultiplier := int(math.Ceil(math.Log2(float64(opts.LengthHint) + 2.0)))

	// Determine a scaling factor to adjust the buffer size.
	// This factor ensures the buffer is sufficiently large to accommodate the randomness needed,
	// balancing between performance (less frequent random reads) and memory usage.
	scalingFactor := int(math.Max(3.0, float64(alphabetLen)/math.Pow(float64(opts.LengthHint), 0.6)))

	// Compute the buffer multiplier by adding the base multiplier and a fraction of the scaling factor.
	// This combination fine-tunes the buffer size, considering both the ID length and the alphabet size.
	bufferMultiplier := baseMultiplier + int(math.Ceil(float64(scalingFactor)/1.5))

	// Calculate the total buffer size in bytes for generating random data.
	// The buffer size is influenced by the buffer multiplier, bytes needed per character,
	// and a factor that scales with the length hint.
	// A larger buffer reduces the number of calls to the random number generator, improving efficiency.
	bufferSize := bufferMultiplier * int(bytesNeeded) * int(math.Max(1.5, float64(opts.LengthHint)/10.0))

	return &runtimeConfig{
		randReader:       opts.RandReader,
		byteAlphabet:     byteAlphabet,
		runeAlphabet:     alphabetRunes,
		mask:             mask,
		bitsNeeded:       bitsNeeded,
		bytesNeeded:      bytesNeeded,
		bufferSize:       bufferSize,
		bufferMultiplier: bufferMultiplier,
		scalingFactor:    scalingFactor,
		baseMultiplier:   baseMultiplier,
		alphabetLen:      alphabetLen,
		isASCII:          isASCII,
		isPowerOfTwo:     isPowerOfTwo,
		lengthHint:       opts.LengthHint,
		maxBytesPerRune:  maxBytesPerRune,
	}, nil
}

// processRandomBytes extracts and returns an unsigned integer from the given randomBytes slice,
// starting at the specified index 'i'. The size of the returned value is determined by the
// g.config.bytesNeeded field.
//
// Parameters:
//   - randomBytes: A byte slice containing random data.
//   - i: The starting index from which to extract the required bytes from the randomBytes slice.
//
// Returns:
//   - uint: An unsigned integer constructed from the bytes, with a size defined by g.config.bytesNeeded.
//
// Behavior:
//   - If bytesNeeded is 1, a single byte is returned as an unsigned integer.
//   - If bytesNeeded is 2, the function returns a 16-bit unsigned integer (2 bytes) in Big Endian order.
//   - If bytesNeeded is 4, the function returns a 32-bit unsigned integer (4 bytes) in Big Endian order.
//   - For other values of bytesNeeded, it constructs an unsigned integer by shifting and combining each byte.
//
// This function is kept small to encourage inlining by the compiler.
//
//go:inline
func (g *generator) processRandomBytes(randomBytes []byte, i int) uint {
	switch g.config.bytesNeeded {
	case 1:
		return uint(randomBytes[i])
	case 2:
		return uint(binary.BigEndian.Uint16(randomBytes[i : i+2]))
	case 4:
		return uint(binary.BigEndian.Uint32(randomBytes[i : i+4]))
	default:
		var rnd uint
		for j := 0; j < int(g.config.bytesNeeded); j++ {
			rnd = (rnd << 8) | uint(randomBytes[i+j])
		}
		return rnd
	}
}

// New generates a new Nano ID string of the specified length.
//
// It validates the provided length to ensure it is a positive integer.
// Depending on the generator's configuration, it generates the ID using the appropriate method.
//
// Parameters:
//   - length int: The desired number of characters in the generated Nano ID.
//
// Returns:
//   - string: The generated Nano ID.
//   - error: An error object if the generation fails due to invalid input.
//
// Error Conditions:
//   - ErrInvalidLength: Returned if the provided length is less than or equal to zero.
//
// Usage Example:
//
//	id, err := DefaultGenerator.New(21)
//	if err != nil {
//	    // handle error
//	}
//	fmt.Println("Generated ID:", id)
func (g *generator) New(length int) (ID, error) {
	if length <= 0 {
		return EmptyID, ErrInvalidLength
	}

	if g.config.isASCII {
		return g.newASCII(length)
	}
	return g.newUnicode(length)
}

// newASCII generates a new Nano ID using the ASCII alphabet.
func (g *generator) newASCII(length int) (ID, error) {
	randomBytesPtr := g.entropyPool.Get().(*[]byte)
	randomBytes := *randomBytesPtr
	bufferLen := len(randomBytes)

	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	mask := g.config.mask
	bytesNeeded := g.config.bytesNeeded
	isPowerOfTwo := g.config.isPowerOfTwo

	// Use strings.Builder to build the ID efficiently
	var sb strings.Builder
	sb.Grow(length) // Preallocate capacity to minimize allocations

	// Defer returning the randomBytes buffer to the pool
	defer func() {
		g.entropyPool.Put(randomBytesPtr)
	}()

	for attempts := 0; cursor < length && attempts < maxAttempts; attempts++ {
		neededBytes := (length - cursor) * int(bytesNeeded)
		if neededBytes > bufferLen {
			neededBytes = bufferLen
		}

		// Fill the random bytes buffer
		if _, err := g.config.randReader.Read(randomBytes[:neededBytes]); err != nil {
			return EmptyID, err
		}

		// Process each segment of random bytes
		for i := 0; i < neededBytes && cursor < length; i += int(bytesNeeded) {
			rnd := g.processRandomBytes(randomBytes, i)
			rnd &= mask

			if isPowerOfTwo || int(rnd) < int(g.config.alphabetLen) {
				sb.WriteByte(g.config.byteAlphabet[rnd])
				cursor++
			}
		}
	}

	// Check for max attempts
	if cursor < length {
		return EmptyID, ErrExceededMaxAttempts
	}

	return ID(sb.String()), nil
}

// newUnicode generates a new Nano ID using the Unicode alphabet.
func (g *generator) newUnicode(length int) (ID, error) {
	// Retrieve random bytes from the pool
	randomBytesPtr := g.entropyPool.Get().(*[]byte)
	randomBytes := *randomBytesPtr
	bufferLen := len(randomBytes)

	cursor := 0
	maxAttempts := length * maxAttemptsMultiplier
	mask := g.config.mask
	bytesNeeded := g.config.bytesNeeded
	isPowerOfTwo := g.config.isPowerOfTwo

	var sb strings.Builder

	// Use the precomputed maximum bytes per rune to estimate buffer size
	sb.Grow(length * g.config.maxBytesPerRune)

	// Defer returning the randomBytes buffer to the pool
	defer func() {
		g.entropyPool.Put(randomBytesPtr)
	}()

	for attempts := 0; cursor < length && attempts < maxAttempts; attempts++ {
		neededBytes := (length - cursor) * int(bytesNeeded)
		if neededBytes > bufferLen {
			neededBytes = bufferLen
		}

		// Fill the random bytes buffer
		if _, err := g.config.randReader.Read(randomBytes[:neededBytes]); err != nil {
			return EmptyID, err
		}

		// Process each segment of random bytes
		for i := 0; i < neededBytes && cursor < length; i += int(bytesNeeded) {
			rnd := g.processRandomBytes(randomBytes, i)
			rnd &= mask

			if isPowerOfTwo || int(rnd) < int(g.config.alphabetLen) {
				sb.WriteRune(g.config.runeAlphabet[rnd])
				cursor++
			}
		}
	}

	// Check for max attempts
	if cursor < length {
		return EmptyID, ErrExceededMaxAttempts
	}

	return ID(sb.String()), nil
}

// Reader is the interface that wraps the basic Read method.
//
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
//
// If len(p) == 0, Read should always return n == 0. It may return a
// non-nil error if some error condition is known, such as EOF.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
func (g *generator) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	length := len(p)
	id, err := g.New(length)
	if err != nil {
		return 0, err
	}

	copy(p, id)
	return length, nil
}

// IsEmpty returns true if the ID is an empty ID (EmptyID)
func (id ID) IsEmpty() bool {
	return id.Compare(EmptyID) == 0
}

// Compare compares two IDs lexicographically and returns an integer.
// The result will be 0 if id==other, -1 if id < other, and +1 if id > other.
//
// Parameters:
//   - other ID: The ID to compare against.
//
// Returns:
//   - int: An integer indicating the comparison result.
//
// Usage:
//
//	id1 := ID("V1StGXR8_Z5jdHi6B-myT")
//	id2 := ID("V1StGXR8_Z5jdHi6B-myT")
//	result := id1.Compare(id2)
//	fmt.Println(result) // Output: 0
func (id ID) Compare(other ID) int {
	return strings.Compare(string(id), string(other))
}

// String returns the string representation of the ID.
// It implements the fmt.Stringer interface, allowing the ID to be
// used seamlessly with fmt package functions like fmt.Println and fmt.Printf.
//
// Example:
//
//	id := Must()
//	fmt.Println(id) // Output: V1StGXR8_Z5jdHi6B-myT
func (id ID) String() string {
	return string(id)
}

// MarshalText converts the ID to a byte slice.
// It implements the encoding.TextMarshaler interface, enabling the ID
// to be marshaled into text-based formats such as XML and YAML.
//
// Returns:
//   - A byte slice containing the ID.
//   - An error if the marshaling fails.
//
// Example:
//
//	id := Must()
//	text, err := id.MarshalText()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(string(text)) // Output: V1StGXR8_Z5jdHi6B-myT
func (id ID) MarshalText() ([]byte, error) {
	return []byte(id), nil
}

// UnmarshalText parses a byte slice and assigns the result to the ID.
// It implements the encoding.TextUnmarshaler interface, allowing the ID
// to be unmarshaled from text-based formats.
//
// Parameters:
//   - text: A byte slice containing the ID data.
//
// Returns:
//   - An error if the unmarshaling fails.
//
// Example:
//
//	var id ID
//	err := id.UnmarshalText([]byte("new-id"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(id) // Output: new-id
func (id *ID) UnmarshalText(text []byte) error {
	*id = ID(text)
	return nil
}

// MarshalBinary converts the ID to a byte slice.
// It implements the encoding.BinaryMarshaler interface, enabling the ID
// to be marshaled into binary formats for efficient storage or transmission.
//
// Returns:
//   - A byte slice containing the ID.
//   - An error if the marshaling fails.
//
// Example:
//
//	id := Must()
//	binaryData, err := id.MarshalBinary()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(binaryData) // Output: [86 49 83 116 71 88 82 56 95 90 ...]
func (id ID) MarshalBinary() ([]byte, error) {
	return []byte(id), nil
}

// UnmarshalBinary parses a byte slice and assigns the result to the ID.
// It implements the encoding.BinaryUnmarshaler interface, allowing the ID
// to be unmarshaled from binary formats.
//
// Parameters:
//   - data: A byte slice containing the binary ID data.
//
// Returns:
//   - An error if the unmarshaling fails.
//
// Example:
//
//	var id ID
//	err := id.UnmarshalBinary([]byte{86, 49, 83, 116, 71, 88, 82, 56, 95, 90}) // "V1StGXR8_Z5jdHi6B-myT"
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(id) // Output: V1StGXR8_Z5jdHi6B-myT
func (id *ID) UnmarshalBinary(data []byte) error {
	*id = ID(data)
	return nil
}

// Config holds the runtime configuration for the Nano ID generator.
//
// It is immutable after initialization and provides all the necessary
// parameters for generating unique IDs efficiently and securely.
func (g *generator) Config() Config {
	return g.config
}

// AlphabetLen returns the number of unique characters in the provided alphabet.
//
// This length determines the range of indices for selecting characters during ID generation.
// Using uint16 allows for alphabets up to 65,535 characters.
func (r runtimeConfig) AlphabetLen() uint16 {
	return r.alphabetLen
}

// BaseMultiplier returns the foundational multiplier used in buffer size calculations.
//
// It is based on the logarithm of the intended ID length (LengthHint) plus 2.
// This helps scale the buffer size appropriately with different ID lengths.
func (r runtimeConfig) BaseMultiplier() int {
	return r.baseMultiplier
}

// BitsNeeded returns the minimum number of bits required to represent all possible indices of the alphabet.
//
// This value is crucial for generating random numbers that map uniformly to the alphabet indices without bias.
func (r runtimeConfig) BitsNeeded() uint {
	return r.bitsNeeded
}

// BufferMultiplier returns the combined multiplier used in the buffer size calculation.
//
// It adds a fraction of the scaling factor to the base multiplier to fine-tune the buffer size,
// considering both the ID length and the alphabet size.
func (r runtimeConfig) BufferMultiplier() int {
	return r.bufferMultiplier
}

// BufferSize returns the total size of the buffer (in bytes) used for generating random data.
//
// The buffer size is calculated to balance efficiency and performance,
// minimizing calls to the random number generator by reading larger chunks of random data at once.
func (r runtimeConfig) BufferSize() int {
	return r.bufferSize
}

// ByteAlphabet returns the slice of bytes representing the alphabet,
// used when the alphabet consists solely of ASCII characters.
//
// For non-ASCII alphabets, this returns nil, and RuneAlphabet is used instead.
func (r runtimeConfig) ByteAlphabet() []byte {
	return r.byteAlphabet
}

// BytesNeeded returns the number of bytes required to store the BitsNeeded for each character in the ID.
//
// It rounds up BitsNeeded to the nearest byte, ensuring sufficient space for random data generation.
func (r runtimeConfig) BytesNeeded() uint {
	return r.bytesNeeded
}

// IsASCII returns true if the alphabet consists solely of ASCII characters.
//
// This allows for optimization in processing, using bytes instead of runes for ID generation.
func (r runtimeConfig) IsASCII() bool {
	return r.isASCII
}

// IsPowerOfTwo returns true if the length of the alphabet is a power of two.
//
// When true, random index selection can be optimized using bitwise operations,
// such as bitwise AND with the mask, improving performance.
func (r runtimeConfig) IsPowerOfTwo() bool {
	return r.isPowerOfTwo
}

// LengthHint returns the intended length of the IDs to be generated.
//
// This hint is used in calculations to adjust buffer sizes and scaling factors accordingly.
func (r runtimeConfig) LengthHint() uint16 {
	return r.lengthHint
}

// Mask returns the bitmask used to extract the necessary bits from randomly generated bytes.
//
// The mask is essential for efficiently mapping random values to valid alphabet indices,
// ensuring uniform distribution and preventing bias.
func (r runtimeConfig) Mask() uint {
	return r.mask
}

// RandReader returns the source of randomness used for generating IDs.
//
// It is typically a cryptographically secure random number generator (e.g., crypto/rand.Reader).
func (r runtimeConfig) RandReader() io.Reader {
	return r.randReader
}

// RuneAlphabet returns the slice of runes representing the alphabet.
//
// This is used for ID generation when the alphabet includes non-ASCII (multibyte) characters,
// allowing support for a wider range of characters.
func (r runtimeConfig) RuneAlphabet() []rune {
	return r.runeAlphabet
}

// ScalingFactor returns the scaling factor used to adjust the buffer size.
//
// It balances the influence of the alphabet size and the intended ID length,
// ensuring efficient random data generation without excessive memory usage.
func (r runtimeConfig) ScalingFactor() int {
	return r.scalingFactor
}

// MaxBytesPerRune represents the maximum number of bytes required to encode
// any rune in the alphabet using UTF-8 encoding.
func (r runtimeConfig) MaxBytesPerRune() int {
	return r.maxBytesPerRune
}
