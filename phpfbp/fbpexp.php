<?php

// -------------------------------------------------------------------------------------
//  What to run, when this file is executed as a script
// -------------------------------------------------------------------------------------

function main() {
	// $fooProgram = new FooExample();
	// $fooProgram->run();
	$splitExample = new StringSplitExample();
	$splitExample->run();
}

// -------------------------------------------------------------------------------------
//  Program Examples
// -------------------------------------------------------------------------------------

class FooExample implements IProgram {
	function run() {
		// Initialize processes
		$foo_writer = new FooWriter();
		$foo_to_bar = new Foo2Bar();
		$printer = new Printer();

		// Connect network
		$foo_to_bar->in_foo = $foo_writer->out_foo;
		$printer->in_str = $foo_to_bar->out_bar;

		$net = new Network();
		$net->add_processes([$foo_writer, $foo_to_bar, $printer]);
		$net->run();
	}
}

class StringSplitExample implements IProgram {
	function run() {
		$str_writer = new StringWriter('hejsan hoppsan hoppas doppas kloppas');
		$str_splitter = new StringSplitter();
		$str_printer = new Printer();

		$str_splitter->in_string = $str_writer->out_string;
		$str_printer->in_str = $str_splitter->out_string_parts;

		$net = new Network();
		$net->add_processes([$str_writer, $str_splitter, $str_printer]);
		$net->run();
	}
}

// -------------------------------------------------------------------------------------
//  Processes
// -------------------------------------------------------------------------------------

class FooWriter implements IProcess {
	public $out_foo = null;

	public function __construct() {
		$this->out_foo = new Chan();
	}

	public function execute() {
		$this->out_foo->send('foo');
		$this->out_foo->close();

		return true;
	}
}

class Foo2Bar implements IProcess {
	public $in_foo = null;
	public $out_bar = null;

	public function __construct() {
		$this->in_foo = new Chan();
		$this->out_bar= new Chan();
	}

	public function execute() {
		$instr = $this->in_foo->recv();
		$outstr = str_replace('foo', 'bar', $instr);
		$this->out_bar->send($outstr);

		if ( $this->in_foo->done() ) {
			$this->out_bar->close();
			return true;
		}
		return false;
	}
}

class StringWriter implements IProcess {
	public $out_string = null;

	protected $string_to_write = '';

	public function __construct($string_to_write) {
		$this->out_string = new Chan();
		$this->string_to_write = $string_to_write;
	}

	public function execute() {
		$this->out_string->send($this->string_to_write);
		$this->out_string->close();
		return true; // Done = true
	}
}

class StringSplitter implements IProcess {
	public $in_string = null;
	public $out_string_parts = null;

	public function __construct() {
		$this->in_string = new Chan();
		$this->out_string_parts = new Chan();
	}

	public function execute() {
		$instr = $this->in_string->recv();
		$string_parts = explode(' ', $instr);
		foreach ($string_parts as $str_part ) {
			$this->out_string_parts->send($str_part);
		}
		$this->out_string_parts->close();
		return true; // Done = true
	}
}

class Printer implements IProcess {
	public $in_str = null;

	public function __construct() {
		$this->in_str = new Chan();
	}

	public function execute() {
		$instr = $this->in_str->recv();
		if ( $instr !== '' ) {
			echo "$instr\n";
		}

		if ( $this->in_str->done() ) {
			return true;
		}
		return false;
	}
}

// -------------------------------------------------------------------------------------
//  FBP Components
// -------------------------------------------------------------------------------------

class Network {
	protected $processes = [];

	public function add_processes($processes) {
		$this->processes = array_merge($this->processes, $processes);
	}

	public function run() {
		while ( count( $this->processes ) > 0 ) {
			foreach ( $this->processes as $i => $p ) {
				$done = $p->execute();
				if ( $done ) {
					unset($this->processes[$i]);
				}
			}
		}
	}
}

class Chan {
	protected $items = array();
	protected $closed = false;

	public function send($item) {
		array_unshift($this->items, $item);
	}

	public function recv() {
		return array_pop($this->items);
	}

	public function close() {
		$this->closed = true;
	}

	public function done() {
		return $this->closed && count( $this->items ) == 0;
	}
}

interface INetwork {
	function run();
}

interface IProgram {
	function run();
}

interface IProcess {
	function execute();
}

// -------------------------------------------------------------------------------------
//  Run this file as a script
// -------------------------------------------------------------------------------------

main();