use curta_core::{utils, CurtaProver, CurtaStdin, CurtaVerifier};

/// The ELF we want to execute inside the zkVM.
const REGEX_IO_ELF: &[u8] =
    include_bytes!("../../../programs/demo/regex/elf/riscv32im-succinct-zkvm-elf");

fn main() {
    // Setup a tracer for logging.
    utils::setup_tracer();

    // Create a new stdin with d the input for the program.
    let mut stdin = CurtaStdin::new();

    let pattern = "a+".to_string();
    let target_string = "an era of truth, not trust".to_string();

    // Write in a simple regex pattern.
    stdin.write(&pattern);
    stdin.write(&target_string);

    // Generate the proof for the given program and input.
    let mut proof = CurtaProver::prove(REGEX_IO_ELF, stdin).expect("proving failed");

    // Read the output.
    let res = proof.stdout.read::<bool>();
    println!("res: {}", res);

    // Verify proof.
    CurtaVerifier::verify(REGEX_IO_ELF, &proof).expect("verification failed");

    // Save the proof.
    proof
        .save("proof-with-pis.json")
        .expect("saving proof failed");

    println!("succesfully generated and verified proof for the program!")
}
