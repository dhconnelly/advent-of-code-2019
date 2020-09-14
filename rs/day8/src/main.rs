const WIDTH: usize = 25;
const HEIGHT: usize = 6;
const LAYER_SIZE: usize = WIDTH * HEIGHT;

#[derive(Debug)]
struct Layer(Vec<char>);

#[derive(Debug)]
struct Image(Vec<Layer>);

fn read_layer(text: &str) -> Layer {
    Layer(text.chars().collect())
}

fn read_image(text: &str) -> Image {
    let n = text.len() / LAYER_SIZE;
    let mut layers = Vec::new();
    for i in 0..n {
        let from = i * LAYER_SIZE;
        let to = from + LAYER_SIZE;
        let layer = read_layer(&text[from..to]);
        layers.push(layer);
    }
    Image(layers)
}

fn count_elem(chs: &[char], elem: char) -> usize {
    chs.iter().filter(|ch| **ch == elem).count()
}

fn find_layer_fewest_zeroes(img: &Image) -> &Layer {
    let zeroes = |chs| count_elem(chs, '0');
    img.0
        .iter()
        .min_by(|Layer(l1), Layer(l2)| zeroes(&l1).cmp(&zeroes(&l2)))
        .unwrap()
}

fn num1s_times_num2s(layer: &Layer) -> usize {
    let ones = |chs| count_elem(chs, '1');
    let twos = |chs| count_elem(chs, '2');
    ones(&layer.0) * twos(&layer.0)
}

#[derive(Debug, Clone, Copy, PartialEq)]
enum Color {
    Black,
    White,
    Transparent,
}

impl Color {
    fn of(ch: char) -> Color {
        match ch {
            '0' => Color::Black,
            '1' => Color::White,
            '2' => Color::Transparent,
            ch => panic!("bad color: {}", ch),
        }
    }
}

#[derive(Debug)]
struct RenderedImage(Vec<Color>);

impl std::fmt::Display for RenderedImage {
    fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
        for i in 0..HEIGHT {
            for j in 0..WIDTH {
                let idx = i * WIDTH + j;
                match self.0[idx] {
                    Color::Black => write!(f, "⬛")?,
                    Color::White => write!(f, "⬜")?,
                    Color::Transparent => write!(f, " ")?,
                }
            }
            println!();
        }
        Ok(())
    }
}

fn apply_layer(ri: &mut RenderedImage, layer: &Layer) {
    assert!(ri.0.len() == layer.0.len());
    for i in 0..ri.0.len() {
        let layer_color = Color::of(layer.0[i]);
        if layer_color != Color::Transparent {
            ri.0[i] = layer_color;
        }
    }
}

fn render(img: &Image) -> RenderedImage {
    let mut ri = RenderedImage(vec![Color::Transparent; LAYER_SIZE]);
    for layer in img.0.iter().rev() {
        apply_layer(&mut ri, layer);
    }
    ri
}

fn main() {
    let path: String = std::env::args().nth(1).unwrap();
    let text: String = std::fs::read_to_string(&path).unwrap();
    let img = read_image(&text);

    println!("{}", num1s_times_num2s(find_layer_fewest_zeroes(&img)));
    println!("{}", render(&img));
}
