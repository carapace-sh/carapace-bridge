use bpaf::Bpaf;

#[derive(Debug, Clone, Bpaf)]
#[bpaf(options, version)]
enum Demo {
    /// Say hello
    #[bpaf(command)]
    Hello {
        /// Name of the person to greet
        name: String,
        /// Be formal
        #[bpaf(short, long)]
        formal: bool,
    },
    /// Create something
    #[bpaf(command)]
    Create {
        /// Name of the thing
        name: String,
        /// Type of thing to create
        #[bpaf(long, argument("KIND"))]
        kind: String,
    },
}

fn main() {
    println!("{:#?}", demo().run());
}
