:orphan:

.. _torch_batch_process_api_ref:

###################################################
 ``name of det torch batch process`` API Reference
###################################################

.. meta::
   :description: Familiarize yourself with the Torch Batch Process API.

+--------------------------------------------+
| User Guide                                 |
+============================================+
| :ref:`torch_batch_processing_ug`           |
+--------------------------------------------+

.. caution::

   This is an experimental API and may change at any time.

The main arguments to torch_batch_process is processor class and dataset.

.. code:: python

   torch_batch_process(
       batch_processor_cls=MyProcessor
       dataset=dataset
   )

**************************************************
 ``determined.pytorch.TorchBatchProcessorContext``
**************************************************

.. autoclass:: determined.pytorch.experimental.TorchBatchProcessorContext
   :members:
   :member-order: bysource


*******************************************
 ``determined.pytorch.TorchBatchProcessor``
*******************************************

.. autoclass:: determined.pytorch.experimental.TorchBatchProcessor
   :members:
   :member-order: bysource

