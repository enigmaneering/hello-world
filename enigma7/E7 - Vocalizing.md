# `E7 - Vocalizing`
### `Alex Petz, Ignite Laboratories, October 2025`

---

### Particles and Air Pressure

The next phase of this is odd - I'm exploring it long before I have the neural code necessary to execute on it, yet
it's critical for me to document my findings as I explore.  The concept of speech reproduction is one I've considered
since I was an infant, but not until recently did I consider the _implementation._  For me, I've envisioned the notion
that sound is a byproduct of the environment _first_ and the instrument _second._  In fact, the instrument was created
as a way of _harnessing_ an environmental quality into a reproducible set of noises.  Drums, for instance, are a great
example of this: the majority of components are simply _cylinders_ with a film stretched over at least one end.  Humanity,
however, _created_ what we call a 'drum' - by recognizing how _nature_ had presented drum-like surfaces which we could
_replicate._

    tl;dr - a drum needs only a taught surface and a stick

So, how does that translate into _code?_  Well, we can digitally model a _drum_ and the mechanics of _air pressure_
affected by vibration.  Up to this point the concept of _neural rendering_ should be pretty well established, so now
we get to consider the ideas of how _particle physics_ can be used to simulate _air pressure._  When I speak of "particles"
my mind instantly connects to _particle emitters_ in code - distinct points that need calculated as they spatially
traverse through a simulated physical space.  The same _absolutely_ applies to vocal reproduction, but to achieve this
we have to be able to _animate_ the surfaces in tandem with our vocal calculations.  We're going to grab into the toolbag
of nearly every topic discussed up to this point, but I promise it's only as complicated as you let yourself believe
it to be!

You've made it this far, brave enigmaneer - you've got this!

### Let's Get Fricative

The ultimate goal of this enigma is _not_ to create a drum, though we'll definitely be starting there - it's to create
_speech!_  Thus, before we even remotely begin, I'd like to take a brief moment to give you a crash course on my naïve
understanding of biological speech production, phonemes, and the International Phonetic Alphabet.  I assure you, I'm 
fully aware I have _zero qualifications_ to be attempting to _teach_ this stuff - please don't take my understanding
to be gospel!  However, my lack of education is no barrier to _my_ attempts at such an absurd challenge =)

So, what even _is_ a phoneme?  Well, if the bare minimum component of outputting a _written_ language is the _lexeme,_
a _phoneme_ is the minimum component of emitting a _spoken_ language.  A phoneme, at its core, is the arrangement of
_articulators_ (such as the tongue, teeth, lips or palate) which, when provided controlled air pressure, will yield
a distinctly useful component of sound.  The International Phonetic Alphabet curates the phonetic categories of
phonemes, of which there are more than a hundred.  Fortunately, _I only speak English!_  Even if I tried to explore
the other phonetic categories not used by the English language, I would have absolutely no awareness if the emitted
sounds were intelligible.  Thus, that limits our phonetic scope to roughly 46 phonemes.  At the time of this writing,
there are _endless_ resources available to understand them - but I'd like to shout out to pronunciationstudio.com,
who've been serving this field since 2008 and have excellent diagrams from which much of my data was sourced.

Using phonemes, and those wonderful cross-section diagrams of the vocal cavities, we can _animate_ the articulators between
positions and apply _air pressure_ behind them.  A phoneme alone is _not_ enough to produce speech, after all!  To generate
speech, we get to _propagate_ a pressure wave through a virtual cavity using particles (and voxels, but we'll get to that
later)

So, the first component of research I'd like to discuss is the eight different English phoneme categories - don't get too
daunted, you don't need deep knowledge to proceed forward

0. Monophthongs - A vowel that has a single perceived auditory quality.
1. Dipthongs - A sound formed by the combination of two vowels in a single syllable, in which the sound begins as one vowel and moves toward another (as in coin, loud, and side)
2. Fricatives - Denoting a type of consonant made by the friction of breath in a narrow opening, producing a turbulent air flow
3. Affricates - A phoneme which combines a plosive with an immediately following fricative sharing the same place of articulation, e.g. ch as in chair and j as in jar
4. Plosives - Denoting a consonant that is produced by stopping the airflow using the lips, teeth, or palate, followed by a sudden release of air
5. Approximates - Sounds that are produced by bringing articulators close together without creating turbulent airflow, unlike a consonant with friction
6. Lateral Approximates - Consonant sounds produced by allowing air to flow freely over the sides of the tongue while a blockage is made in the middle of the mouth
7. Nasals - Sounds produced when air is directed through the nose, as the soft palate (velum) is lowered to allow air to escape from the nasal cavity

My distillation of these definitions is that they define the quality of how _air_ is involved in the process of speech
production.  Should pressure be built and then released once the articulators are in position, like with the nasals?
Or should it be given a turbulent flow by blocking with the teeth, as in approximations?  In addition to this quality
is the _dipthong,_ which implies the articulators must move _during_ the phoneme's pronunciation, rather than simply
before processing that phoneme.  All of these qualities are important in speech production, but neatly pack into the
idea of _particles_ propagating a pressure wave through a virtual cavity.

To produce speech, we have several critical biological components as well

0. The Tongue - Responsible for the changing the path of the air
1. The Teeth - Responsible for creating turbulent air
2. The Palate - Provides the shape of the cavity
3. The Larynx - Controls the aperture of air
4. The Vocal Folds - Act as the oscillatory pressure source
5. The Soft Palate - Switches between nasal and oral sound emission

Lucky for us, an instrument such as the voicebox can be considered analogous to any other musical instrument - albeit
_immensely_ more nuanced between one another.  If you were to stretch a trumpet out as a single long tube, with the
valves still placed at their appropriate points, the instrument would sound theoretically identical to the more compact
form we know and love today.  We can leverage this quality ourselves in the creation of our vocal _instrument!_  Our
process will start by evolving the idea of a straight piped larynx controlling a tongue emitter as if the vocal folds
were a reed on a saxophone (how did I wind up writing this stuff!?)

So - let's begin!