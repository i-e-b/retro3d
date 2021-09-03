package pkg

type mappingVars struct {
	dizdx, duizdx, dvizdx, dizdy, duizdy, dvizdy float64
	dizdxn, duizdxn, dvizdxn                     float64
	xa, xb, iza, uiza, viza                      float64
	dxdya, dxdyb, dizdya, duizdya, dvizdya       float64

	a, b, c      Vec3t
	tex          *Texture
	screen       *[]uint32
	screenWidth  int
	screenHeight int
}

// fixed point stuff
const subDivShift int = 4
const subDivSize int64 = 1 << subDivShift

func TextureTriangle(a Vec3t, b Vec3t, c Vec3t, tex *Texture, buf *[]uint32, width int, height int) {
	poly := mappingVars{
		a: a, b: b, c: c,
		tex: tex, screen: buf,
		screenWidth: width, screenHeight: height,
	}
	drawTexPolyPerspDivSubTri(&poly)
}

func drawTexPolyPerspDivSubTri(poly *mappingVars) {
	var x1, y1, x2, y2, x3, y3 float64
	var iz1, uiz1, viz1, iz2, uiz2, viz2, iz3, uiz3, viz3 float64
	var dxdy1 = 0.0
	var dxdy2 = 0.0
	var dxdy3 = 0.0
	var denom = 0.0
	var dy float64
	var y1i, y2i, y3i int

	// Shift XY coordinate system (+0.5, +0.5) to match the subpixel strategy technique
	x1 = poly.a.X + 0.5
	y1 = poly.a.Y + 0.5
	x2 = poly.b.X + 0.5
	y2 = poly.b.Y + 0.5
	x3 = poly.c.X + 0.5
	y3 = poly.c.Y + 0.5

	// Calculate alternative 1/Z, U/Z and V/Z values which will be interpolated
	iz1 = 1.0 / poly.a.Z
	iz2 = 1.0 / poly.b.Z
	iz3 = 1.0 / poly.c.Z
	uiz1 = float64(poly.a.U) * iz1
	viz1 = float64(poly.a.V) * iz1
	uiz2 = float64(poly.b.U) * iz2
	viz2 = float64(poly.b.V) * iz2
	uiz3 = float64(poly.c.U) * iz3
	viz3 = float64(poly.c.V) * iz3

	// Sort the vertices in increasing Y order
	if y1 > y2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		iz1, iz2 = iz2, iz1
		uiz1, uiz2 = uiz2, uiz1
		viz1, viz2 = viz2, viz1
	}

	if y1 > y3 {
		x1, x3 = x3, x1
		y1, y3 = y3, y1
		iz1, iz3 = iz3, iz1
		uiz1, uiz3 = uiz3, uiz1
		viz1, viz3 = viz3, viz1
	}

	if y2 > y3 {
		x2, x3 = x3, x2
		y2, y3 = y3, y2
		iz2, iz3 = iz3, iz2
		uiz2, uiz3 = uiz3, uiz2
		viz2, viz3 = viz3, viz2
	}

	y1i = int(y1)
	y2i = int(y2)
	y3i = int(y3)

	// Skip poly if it's too thin to cover any pixels at all
	if y1i == y2i && y1i == y3i || (int(x1) == int(x2) && int(x1) == int(x3)) {
		return
	}

	// Calculate horizontal and vertical increments for UV axes (these
	//  calculations are certainly not optimal, although they're stable
	//  (handles any dy being 0)

	denom = (x3-x1)*(y2-y1) - (x2-x1)*(y3-y1)

	// Skip poly if it's an infinitely thin line
	if denom == 0 {
		return
	}

	denom = 1.0 / denom // Reciprocal for speeding up
	poly.dizdx = ((iz3-iz1)*(y2-y1) - (iz2-iz1)*(y3-y1)) * denom
	poly.duizdx = ((uiz3-uiz1)*(y2-y1) - (uiz2-uiz1)*(y3-y1)) * denom
	poly.dvizdx = ((viz3-viz1)*(y2-y1) - (viz2-viz1)*(y3-y1)) * denom
	poly.dizdy = ((iz2-iz1)*(x3-x1) - (iz3-iz1)*(x2-x1)) * denom
	poly.duizdy = ((uiz2-uiz1)*(x3-x1) - (uiz3-uiz1)*(x2-x1)) * denom
	poly.dvizdy = ((viz2-viz1)*(x3-x1) - (viz3-viz1)*(x2-x1)) * denom

	// Horizontal increases for 1/Z, U/Z and V/Z which step one full span
	//  ahead

	poly.dizdxn = poly.dizdx * float64(subDivSize)
	poly.duizdxn = poly.duizdx * float64(subDivSize)
	poly.dvizdxn = poly.dvizdx * float64(subDivSize)

	// Calculate X-slopes along the edges
	if y2 > y1 {
		dxdy1 = (x2 - x1) / (y2 - y1)
	}
	if y3 > y1 {
		dxdy2 = (x3 - x1) / (y3 - y1)
	}
	if y3 > y2 {
		dxdy3 = (x3 - x2) / (y3 - y2)
	}

	// Determine which side of the poly the longer edge is on
	var side = dxdy2 > dxdy1

	if int(y1) == int(y2) {
		side = x1 > x2
	}
	if int(y2) == int(y3) {
		side = x3 > x2
	}

	if !side {
		// Longer edge is on the left side
		// Calculate slopes along left edge

		poly.dxdya = dxdy2
		poly.dizdya = dxdy2*poly.dizdx + poly.dizdy
		poly.duizdya = dxdy2*poly.duizdx + poly.duizdy
		poly.dvizdya = dxdy2*poly.dvizdx + poly.dvizdy

		// Perform subpixel pre-stepping along left edge

		dy = 1 - (y1 - float64(y1i))
		poly.xa = x1 + dy*poly.dxdya
		poly.iza = iz1 + dy*poly.dizdya
		poly.uiza = uiz1 + dy*poly.duizdya
		poly.viza = viz1 + dy*poly.dvizdya

		if y1i < y2i { // Draw upper segment if possibly visible

			// Set right edge X-slope and perform subpixel pre-
			//  stepping
			poly.xb = x1 + dy*dxdy1
			poly.dxdyb = dxdy1

			drawTexPolyPerspDivSubTriSegment(poly, y1i, y2i)
		}

		if y2i < y3i { // Draw lower segment if possibly visible

			// Set right edge X-slope and perform subpixel pre-
			//  stepping

			poly.xb = x2 + (1-(y2-float64(y2i)))*dxdy3
			poly.dxdyb = dxdy3

			drawTexPolyPerspDivSubTriSegment(poly, y2i, y3i)
		}
	} else { // Longer edge is on the right side

		// Set right edge X-slope and perform subpixel pre-stepping

		poly.dxdyb = dxdy2
		dy = 1 - (y1 - float64(y1i))
		poly.xb = x1 + dy*poly.dxdyb

		if y1i < y2i { // Draw upper segment if possibly visible

			// Set slopes along left edge and perform subpixel
			//  pre-stepping

			poly.dxdya = dxdy1
			poly.dizdya = dxdy1*poly.dizdx + poly.dizdy
			poly.duizdya = dxdy1*poly.duizdx + poly.duizdy
			poly.dvizdya = dxdy1*poly.dvizdx + poly.dvizdy
			poly.xa = x1 + dy*poly.dxdya
			poly.iza = iz1 + dy*poly.dizdya
			poly.uiza = uiz1 + dy*poly.duizdya
			poly.viza = viz1 + dy*poly.dvizdya

			drawTexPolyPerspDivSubTriSegment(poly, y1i, y2i)
		}

		if y2i < y3i { // Draw lower segment if possibly visible

			// Set slopes along left edge and perform subpixel
			//  pre-stepping

			poly.dxdya = dxdy3
			poly.dizdya = dxdy3*poly.dizdx + poly.dizdy
			poly.duizdya = dxdy3*poly.duizdx + poly.duizdy
			poly.dvizdya = dxdy3*poly.dvizdx + poly.dvizdy
			dy = 1 - (y2 - float64(y2i))
			poly.xa = x2 + dy*poly.dxdya
			poly.iza = iz2 + dy*poly.dizdya
			poly.uiza = uiz2 + dy*poly.duizdya
			poly.viza = viz2 + dy*poly.dvizdya

			drawTexPolyPerspDivSubTriSegment(poly, y2i, y3i)
		}
	}
}

func drawTexPolyPerspDivSubTriSegment(poly *mappingVars, y1, y2 int) {
	// Bounds checks

	if y1 < 0 {y1 = 0} else if y1 >= poly.screenHeight{y1 = poly.screenHeight-1}
	if y2 < 0 {y2 = 0} else if y2 >= poly.screenHeight{y2 = poly.screenHeight-1}

	texW := int64(poly.tex.Width - 1)
	texH := int64(poly.tex.Height - 1)

	for y1 < y2 { // Loop through all lines in segment
		if y1 >= 0 {
			var iz, uiz, viz float64
			var u1, v1, u2, v2, u, v, du, dv int64
			var x1 = int(poly.xa)
			var x2 = int(poly.xb)

			if x1 < 0 {
				x1 = 0
			} else if x1 >= poly.screenWidth {
				x1 = poly.screenWidth - 1
			}
			if x2 < 0 {
				x2 = 0
			} else if x2 >= poly.screenWidth {
				x2 = poly.screenWidth - 1
			}

			// Perform sub-texel pre-stepping on 1/Z, U/Z and V/Z

			var dx = 1 - (poly.xa - float64(x1))
			iz = poly.iza + dx*poly.dizdx
			uiz = poly.uiza + dx*poly.duizdx
			viz = poly.viza + dx*poly.dvizdx

			var cursor = y1*poly.screenWidth + x1 // for poly.screen

			// Calculate UV for the first pixel

			var z = 65536 / iz
			u2 = int64(uiz * z)
			v2 = int64(viz * z)

			// Length of line segment
			var xcount = int64(x2 - x1)

			for xcount >= subDivSize { // Draw all full-length

				//  spans
				// Step 1/Z, U/Z and V/Z to the next span

				iz += poly.dizdxn
				uiz += poly.duizdxn
				viz += poly.dvizdxn

				u1 = u2
				v1 = v2

				// Calculate UV at the beginning of next span

				z = 65536 / iz
				u2 = int64(uiz * z)
				v2 = int64(viz * z)

				u = u1
				v = v1

				// Calculate linear UV slope over span

				du = (u2 - u1) >> subDivShift
				dv = (v2 - v1) >> subDivShift

				// do jitter smoothing if texels cover more than one pixel
				var qq int64 = 1 << 15
				if (du>>15) > 1 || (dv>>15) > 1 {
					qq = 0
				}
				var j = qq * int64((x1+y1)%2)

				var x = subDivSize
				for x > 0 { // Draw span
					x--
					// Copy pixel from texture to screen

					// jitter the u and v sample points by half a texel
					// on alternating output pixels as a crude anti-alias
					j = qq - j
					var ju = u + j
					var jv = v + j

					var intU = ju >> 16
					var intV = jv >> 16
					if intU > texW {
						intU = texW
					} else if intU < 0 {
						intU = 0
					}
					if intV > texH {
						intV = texH
					} else if intV < 0 {
						intV = 0
					}

					color := poly.tex.Bmp[intV][intU]
					(*poly.screen)[cursor] = color
					cursor++

					// Step horizontally along UV axes

					u += du
					v += dv
				}

				xcount -= subDivSize // One span less
			}

			if xcount != 0 { // Draw last, non-full-length span

				// Step 1/Z, U/Z and V/Z to end of span

				iz += poly.dizdx * float64(xcount)
				uiz += poly.duizdx * float64(xcount)
				viz += poly.dvizdx * float64(xcount)

				u1 = u2
				v1 = v2

				// Calculate UV at end of span

				z = 65536 / iz
				u2 = int64(uiz * z)
				v2 = int64(viz * z)

				u = u1
				v = v1

				// Calculate linear UV slope over span

				du = (u2 - u1) / xcount
				dv = (v2 - v1) / xcount

				// do jitter smoothing if texels cover more than one pixel
				var qq int64 = 1 << 15
				if (du>>15) > 1 || (dv>>15) > 1 {
					qq = 0
				}
				var j = qq * int64((x1+y1)%2)

				for xcount > 0 { // Draw span
					xcount--

					// jitter the u and v sample points by half a texel
					// on alternating output pixels as a crude anti-alias
					j = qq - j
					var ju = u + j
					var jv = v + j

					var intU = ju >> 16
					var intV = jv >> 16
					if intU > texW {
						intU = texW
					} else if intU < 0 {
						intU = 0
					}
					if intV > texH {
						intV = texH
					} else if intV < 0 {
						intV = 0
					}

					// Copy pixel from texture to screen
					color := poly.tex.Bmp[intV][intU]
					(*poly.screen)[cursor] = color
					cursor++

					// Step horizontally along UV axes
					u += du
					v += dv
				}
			}
		}
		// Step vertically along both edges

		poly.xa += poly.dxdya
		poly.xb += poly.dxdyb
		poly.iza += poly.dizdya
		poly.uiza += poly.duizdya
		poly.viza += poly.dvizdya

		y1++
	}
}
